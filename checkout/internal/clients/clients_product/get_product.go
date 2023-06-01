package clients_product

import (
	"context"
	"log"
	"route256/checkout/internal/pkg/product-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	conn, err := grpc.Dial(addressProduct, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := product.NewProductServiceClient(conn)

	return c.GetProduct(ctx, req)
}
