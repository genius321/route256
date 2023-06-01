package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/product"
	"route256/checkout/internal/pkg/checkout"

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
	stocksResp, _ := loms.Stocks(ctx, req.Sku)
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

func (s *service) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	log.Printf("%+v", req)
	type item struct {
		sku   uint32
		count uint32
	}
	items := []item{{sku: 4487693, count: 15}, {sku: 32956725, count: 31}}
	resp := []*checkout.Item{}
	var totalPrice uint32
	for _, v := range items {
		log.Println(v.sku)
		product, err := product.GetProduct(ctx, v.sku)
		if err != nil {
			return nil, fmt.Errorf("get product: %w", err)
		}
		resp = append(resp, &checkout.Item{
			Sku:   v.sku,
			Count: v.count,
			Name:  product.Name,
			Price: product.Price,
		})
		totalPrice += product.Price * v.count
	}
	return &checkout.ListCartResponse{Items: resp, TotalPrice: totalPrice}, nil
}

func (s *service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)
	resp, _ := loms.CreateOrder(ctx, req.User)
	return &checkout.PurchaseResponse{OrderId: resp.OrderId}, nil
}
