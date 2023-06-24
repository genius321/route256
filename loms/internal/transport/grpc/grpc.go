package grpc

import (
	"context"
	"fmt"
	"log"
	orderModels "route256/loms/internal/models/order"
	stockModels "route256/loms/internal/models/stock"
	"route256/loms/internal/pkg/loms"
	"route256/loms/internal/service"
)

type Grpc struct {
	loms.UnimplementedLomsServer
	service *service.Service
}

func NewGrpc(service *service.Service) *Grpc {
	return &Grpc{service: service}
}

func (g *Grpc) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	user := orderModels.User(req.User)
	items := make(orderModels.Items, len(req.Items))
	for i, v := range req.Items {
		items[i].Sku = orderModels.Sku(v.Sku)
		items[i].Count = orderModels.Count(v.Count)
	}

	orderId, err := g.service.CreateOrder(ctx, user, items)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return &loms.CreateOrderResponse{OrderId: int64(orderId)}, nil
}

func (g *Grpc) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	sku := stockModels.Sku(req.Sku)

	stocks, err := g.service.Stocks(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("stocks: %w", err)
	}

	resStocks := make([]*loms.Stock, 0, len(stocks))
	for _, v := range stocks {
		resStocks = append(resStocks, &loms.Stock{
			WarehouseId: int64(v.WarehouseId),
			Count:       uint64(v.Count),
		})
	}

	return &loms.StocksResponse{Stocks: resStocks}, nil
}

func (g *Grpc) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	orderId := orderModels.OrderId(req.OrderId)

	status, user, items, err := g.service.ListOrder(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("listOrder: %w", err)
	}

	resItems := make([]*loms.Item, 0, len(items))
	for _, v := range items {
		resItems = append(resItems, &loms.Item{Sku: uint32(v.Sku), Count: uint32(v.Count)})
	}

	return &loms.ListOrderResponse{
		Status: string(status),
		User:   int64(user),
		Items:  resItems,
	}, nil
}
