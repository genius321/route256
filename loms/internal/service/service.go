package service

import (
	"context"
	"log"
	"route256/loms/internal/pkg/loms"
)

type Service struct {
	loms.UnimplementedLomsServer
}

func (s *Service) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	log.Printf("%+v", req)
	return &loms.CreateOrderResponse{OrderId: 666}, nil
}

func NewLomsServer() *Service {
	return &Service{}
}
