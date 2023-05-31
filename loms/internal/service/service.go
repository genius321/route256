package service

import (
	"context"
	"log"
	"route256/loms/internal/pkg/loms"
)

type service struct {
	loms.UnimplementedLomsServer
}

func (s *service) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	log.Printf("%+v", req)
	return &loms.CreateOrderResponse{OrderId: 666}, nil
}

func NewLomsServer() *service {
	return &service{}
}
