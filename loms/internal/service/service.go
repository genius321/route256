package service

import (
	"context"
	"log"
	"route256/loms/internal/pkg/loms"

	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	loms.UnimplementedLomsServer
}

func NewLomsServer() *service {
	return &service{}
}

func (s *service) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &loms.CreateOrderResponse{OrderId: 666}, nil
}

func (s *service) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &loms.ListOrderResponse{
		Status: "new",
		User:   111,
		Items: []*loms.Item{
			{Sku: 222, Count: 333},
			{Sku: 444, Count: 555},
		},
	}, nil
}
func (s *service) OrderPayed(ctx context.Context, req *loms.OrderPayedRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *service) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *service) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &loms.StocksResponse{
		Stocks: []*loms.Stock{
			{WarehouseId: 1, Count: 150},
			{WarehouseId: 2, Count: 50},
		},
	}, nil
}
