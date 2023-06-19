package main

import (
	"context"
	"fmt"
	"log"
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
	"time"

	_ "net/http/pprof"

	"github.com/aitsvet/debugcharts"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// reset RPS counter
	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			debugcharts.RPS.Set(0)
		}
	}()

	// up debugcharts server
	go func() {
		log.Println("server Pprof on :6060")
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	connToLoms, err := grpc.Dial(config.AddressLoms,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer connToLoms.Close()

	lomsClient := loms.NewLomsClient(connToLoms)

	connToProduct, err := grpc.Dial(config.AddressProduct,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer connToProduct.Close()

	productClient := product.NewProductServiceClient(connToProduct)

	// connection to db
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("CHECKOUT_DATABASE_URL"))
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)
	ratelimiter := ratelimit.New(context.Background(), config.RateLimit)

	// checkout server setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, service.NewCheckoutServer(
		lomsClient,
		productClient,
		repo,
		provider,
		ratelimiter,
	))

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		lis.Addr().String(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	mux := runtime.NewServeMux()
	err = checkout.RegisterCheckoutHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: mux,
	}

	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	log.Fatalln(gwServer.ListenAndServe())
}
