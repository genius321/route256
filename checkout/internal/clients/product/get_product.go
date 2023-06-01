package product

import (
	"context"
	"log"
	"route256/checkout/internal/pkg/product-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addressProduct = "route256.pavl.uk:8082"
	token          = "testtoken"
)

func GetProduct(ctx context.Context, sku uint32) (*product.GetProductResponse, error) {
	conn, err := grpc.Dial(addressProduct, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := product.NewProductServiceClient(conn)

	log.Println(token)
	log.Println(sku)
	product, err := c.GetProduct(ctx, &product.GetProductRequest{
		Token: token,
		Sku:   sku,
	})
	log.Println(product)
	log.Println(err)
	return product, err
}
