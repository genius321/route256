package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"route256/checkout/internal/config"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"
	"route256/checkout/internal/pkg/product-service"
	"route256/checkout/internal/pkg/ratelimit"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/internal/repository/postgres/tx"
	"route256/checkout/internal/service"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracer"
	"time"

	_ "net/http/pprof"

	"github.com/aitsvet/debugcharts"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("ERR loading .env file")
	}
	err = config.Init()
	if err != nil {
		logger.Fatal("ERR config.Init: ", err)
	}
	flag.Parse()
	// Init logger for environment
	logger.SetLoggerByEnvironment(*environment)

	// Init tracer
	if err := tracer.InitGlobal(config.AppConfig.ServiceName); err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// reset RPS counter
	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			debugcharts.RPS.Set(0)
		}
	}()

	// up debugcharts server
	go func() {
		logger.Info("server Pprof on ", fmt.Sprintf(":%d", config.AppConfig.Services.Pprof))
		logger.Info(http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.Services.Pprof), nil))
	}()

	connToLoms, err := grpc.Dial(config.AppConfig.Services.Loms,
		grpc.WithUnaryInterceptor(metrics.MiddlewareClientGRPC),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect to server: %v", err)
	}
	defer connToLoms.Close()

	lomsClient := loms.NewLomsClient(connToLoms)

	connToProduct, err := grpc.Dial(config.AppConfig.Services.ProductService,
		grpc.WithUnaryInterceptor(metrics.MiddlewareClientGRPC),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect to server: %v", err)
	}
	defer connToProduct.Close()

	productClient := product.NewProductServiceClient(connToProduct)

	// connection to db
	pool, err := pgxpool.Connect(ctx, os.Getenv("CHECKOUT_DATABASE_URL"))
	if err != nil {
		logger.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)
	ratelimiter := ratelimit.New(ctx, config.AppConfig.RateLimit)

	// checkout server setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig.GrpcPort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		tracer.MiddlewareGRPC,
		metrics.MiddlewareServerGRPC,
	))
	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, service.NewCheckoutServer(
		lomsClient,
		productClient,
		repo,
		provider,
		ratelimiter,
	))

	logger.Info("server listening at %v", lis.Addr())

	go func() {
		if err = s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		ctx,
		lis.Addr().String(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("Failed to dial server: %v", err)
	}

	mux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "x-trace-id":
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	if err := mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		logger.Fatalln("something wrong with metrics handler", err)
	}

	err = checkout.RegisterCheckoutHandler(ctx, mux, conn)
	if err != nil {
		logger.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppConfig.HttpPort),
		Handler: mux,
	}

	logger.Info("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	logger.Fatal(gwServer.ListenAndServe())
}
