package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"route256/checkout/internal/clients/clients_loms"
	"route256/checkout/internal/clients/clients_product"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"
	"route256/checkout/internal/pkg/product-service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	checkout.UnimplementedCheckoutServer
}

func NewCheckoutServer() *service {
	return &service{}
}

func (s *service) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	stocksResp, _ := clients_loms.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
	stocks := stocksResp.Stocks
	log.Printf("stocks: %v", stocks)
	counter := int64(req.Count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return &emptypb.Empty{}, nil
		}
	}
	return &emptypb.Empty{}, errors.New("stock insufficient")
}

func (s *service) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	return &emptypb.Empty{}, nil
}

func (s *service) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	log.Printf("%+v", req)
	type item struct {
		sku   uint32
		count uint32
	}
	items := []item{{sku: 4487693, count: 15}, {sku: 32956725, count: 31}}
	respItems := make([]*checkout.Item, len(items))
	var totalPrice uint32
	for i, v := range items {
		product, err := clients_product.GetProduct(ctx, &product.GetProductRequest{
			Token: clients_product.Token,
			Sku:   v.sku,
		})
		if err != nil {
			return nil, fmt.Errorf("get product: %w", err)
		}
		respItems[i] = &checkout.Item{
			Sku:   v.sku,
			Count: v.count,
			Name:  product.Name,
			Price: product.Price,
		}

		totalPrice += product.Price * v.count
	}
	return &checkout.ListCartResponse{Items: respItems, TotalPrice: totalPrice}, nil
}

func (s *service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)
	createOrderResp, err := clients_loms.CreateOrder(ctx, &loms.CreateOrderRequest{
		User: req.User,
		Items: []*loms.Item{
			{Sku: 4487693, Count: 15},
			{Sku: 32956725, Count: 31},
		},
	})
	if err != nil {
		return nil, err
	}
	return &checkout.PurchaseResponse{OrderId: createOrderResp.OrderId}, nil
}
