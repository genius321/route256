package service

import (
	"context"
	"log"
	clients "route256/checkout/internal/clients/loms"
	"route256/checkout/internal/pkg/checkout"
)

type service struct {
	checkout.UnimplementedCheckoutServer
}

func NewCheckoutServer() *service {
	return &service{}
}

func (s *service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)
	resp, _ := clients.CreateOrder(ctx, req.User)
	return &checkout.PurchaseResponse{OrderId: resp.OrderId}, nil
}
