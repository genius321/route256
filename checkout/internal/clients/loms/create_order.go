package clients

import (
	"context"
	"log"
	"route256/checkout/internal/pkg/loms"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addressLoms = "localhost:8081"
)

func CreateOrder(ctx context.Context, user int64) (*loms.CreateOrderResponse, error) {
	conn, err := grpc.Dial(addressLoms, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := loms.NewLomsClient(conn)

	return c.CreateOrder(ctx, &loms.CreateOrderRequest{
		User: user,
		Items: []*loms.Item{
			{Sku: 4487693, Count: 15},
			{Sku: 32956725, Count: 31},
		},
	})
}
