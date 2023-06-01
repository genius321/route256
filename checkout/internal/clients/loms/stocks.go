package loms

import (
	"context"
	"log"
	"route256/checkout/internal/pkg/loms"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Stocks(ctx context.Context, sku uint32) (*loms.StocksResponse, error) {
	conn, err := grpc.Dial(addressLoms, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := loms.NewLomsClient(conn)

	return c.Stocks(ctx, &loms.StocksRequest{
		Sku: sku,
	})
}
