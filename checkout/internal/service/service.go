package service

import (
	"context"
	"log"
	"route256/checkout/internal/pkg/checkout"
)

type service struct {
	checkout.UnimplementedCheckoutServer
}

func (s *service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)
	return &checkout.PurchaseResponse{OrderId: 17}, nil
}

func NewCheckoutServer() *service {
	return &service{}
}
