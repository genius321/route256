package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"route256/libs/logger"
	"route256/loms/internal/business"
	"route256/loms/internal/kafka"
	"route256/loms/internal/pkg/loms"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/repository/postgres/tx"
	"route256/loms/internal/service"
	"route256/loms/internal/status"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50052
	httpPort = 8081
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

func main() {
	flag.Parse()

	// Init logger for environment
	logger.SetLoggerByEnvironment(*environment)
	// kafkaProducer
	var brokers = []string{
		"kafka1:29091",
		"kafka2:29092",
		"kafka3:29093",
	}
	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		logger.Fatal(err)
	}
	statusSender := status.NewKafkaSender(kafkaProducer, "statuses")

	// connection to db
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("LOMS_DATABASE_URL"))
	if err != nil {
		logger.Fatal("connect to db: %w", err)
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)

	// loms server setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	loms.RegisterLomsServer(s, service.NewService(business.NewBusiness(repo, provider, statusSender)))

	logger.Info("server listening at %v", lis.Addr())

	go func() {
		if err = s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
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
		logger.Fatalf("Failed to dial server: %v", err)
	}

	mux := runtime.NewServeMux()
	err = loms.RegisterLomsHandler(context.Background(), mux, conn)
	if err != nil {
		logger.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mux,
	}

	logger.Info("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	logger.Fatal(gwServer.ListenAndServe())
}
