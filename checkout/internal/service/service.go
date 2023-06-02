package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"
	"route256/checkout/internal/pkg/product-service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	checkout.UnimplementedCheckoutServer
	lomsClient    loms.LomsClient
	productClient product.ProductServiceClient
}

func NewCheckoutServer(lomsClient loms.LomsClient,
	productClient product.ProductServiceClient) *service {
	return &service{
		lomsClient:    lomsClient,
		productClient: productClient,
	}
}

func (s *service) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	stocksResp, _ := s.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
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
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *service) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	type item struct {
		sku   uint32
		count uint32
	}
	items := []item{{sku: 4487693, count: 15}, {sku: 32956725, count: 31}}
	respItems := make([]*checkout.Item, len(items))
	var totalPrice uint32
	for i, v := range items {
		product, err := s.productClient.GetProduct(ctx, &product.GetProductRequest{
			Token: config.Token,
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
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	createOrderResp, err := s.lomsClient.CreateOrder(ctx, &loms.CreateOrderRequest{
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
