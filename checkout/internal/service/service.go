package service

import (
	"context"
	"log"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct {
	checkout.UnimplementedCheckoutServer
}

func NewCheckoutServer() *service {
	return &service{}
}

const (
	addressLoms = "localhost:8081"
)

func (s *service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)

	conn, err := grpc.Dial(addressLoms, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := loms.NewLomsClient(conn)

	resp, err := c.CreateOrder(ctx, &loms.CreateOrderRequest{
		User: req.User,
		Items: []*loms.Item{
			{Sku: 4487693, Count: 15},
			{Sku: 32956725, Count: 31},
		},
	})
	if err != nil {
		log.Fatalf("Err Purchase: %v", err)
	}

	return &checkout.PurchaseResponse{OrderId: resp.OrderId}, nil
}
