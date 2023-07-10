package service

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/business"
	orderModels "route256/loms/internal/models/order"
	stockModels "route256/loms/internal/models/stock"
	"route256/loms/internal/pkg/loms"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	loms.UnimplementedLomsServer
	service *business.Business
}

func NewService(service *business.Business) *Service {
	return &Service{service: service}
}

func (g *Service) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/service/CreateOrder")
	defer span.Finish()
	logger.Infof("%+v", req)
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

func (g *Service) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/service/Stocks")
	defer span.Finish()
	logger.Infof("%+v", req)
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

func (g *Service) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/service/ListOrder")
	defer span.Finish()
	logger.Infof("%+v", req)
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

func (g *Service) OrderPayed(ctx context.Context, req *loms.OrderPayedRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/service/OrderPayed")
	defer span.Finish()
	logger.Infof("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	err = g.service.OrderPayed(ctx, orderModels.OrderId(req.OrderId))
	if err != nil {
		return nil, fmt.Errorf("orderPayed: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (g *Service) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/service/CancelOrder")
	defer span.Finish()
	logger.Infof("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	err = g.service.CancelOrder(ctx, orderModels.OrderId(req.OrderId))
	if err != nil {
		return nil, fmt.Errorf("cancelOrder: %w", err)
	}
	return &emptypb.Empty{}, nil
}
