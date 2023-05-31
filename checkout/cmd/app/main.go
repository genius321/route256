package main

import (
	"fmt"
	"log"
	"net"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 8080

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, service.NewCheckoutServer())

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
