package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"route256/loms/internal/pkg/loms"
	"route256/loms/internal/repository/schema"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	Serializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Repository interface {
	CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error)
	ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error)
	OrderPayed(ctx context.Context, req *loms.OrderPayedRequest) (*emptypb.Empty, error)
	CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error)

	Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error)
	AddSkuStock(ctx context.Context, sku int64, amount int64, warehouseID int64) (int64, error)
	TakeSkuStock(ctx context.Context, sku int64, amount int64, warehouseID int64) (int64, error)
	AddSkuStockReserve(ctx context.Context, sku int64, amount int64, warehouseID int64, orderId int64) (int64, error)
	DeleteStocksReserveByOrderId(ctx context.Context, orderId int64) error
	TakeStocksReserveByOrderId(ctx context.Context, orderId int64) ([]schema.Stocks, error)
}

type service struct {
	loms.UnimplementedLomsServer
	Repository
	TransactionManager
}

func NewLomsServer(r Repository, t TransactionManager) *service {
	return &service{Repository: r, TransactionManager: t}
}

func (s *service) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	var res *loms.CreateOrderResponse
	err = s.Serializable(ctx, func(ctxTx context.Context) error {
		res, err = s.Repository.CreateOrder(ctxTx, req)
		if err != nil {
			return err
		}
		for _, v := range req.Items {
			warehouseIdReserveCnt, err := s.ReserveStock(ctxTx, int64(v.Sku), int64(v.Count))
			if err != nil {
				return err
			}
			for warehouseID, reserveCnt := range warehouseIdReserveCnt {
				_, err = s.Repository.TakeSkuStock(ctxTx, int64(v.Sku), reserveCnt, warehouseID)
				if err != nil {
					return err
				}
				_, err = s.Repository.AddSkuStockReserve(ctxTx, int64(v.Sku), reserveCnt, warehouseID, res.OrderId)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	return res, nil
}

func (s *service) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return s.Repository.ListOrder(ctx, req)
}

func (s *service) OrderPayed(ctx context.Context, req *loms.OrderPayedRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	var res *emptypb.Empty
	err = s.Serializable(ctx, func(ctxTx context.Context) error {
		err = s.Repository.DeleteStocksReserveByOrderId(ctxTx, req.OrderId)
		if err != nil {
			return err
		}
		res, err = s.Repository.OrderPayed(ctxTx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("order payed: %w", err)
	}
	return res, nil
}

func (s *service) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	var res *emptypb.Empty
	err = s.Serializable(ctx, func(ctxTx context.Context) error {
		reserve, err := s.Repository.TakeStocksReserveByOrderId(ctxTx, req.OrderId)
		if err != nil {
			return err
		}
		for _, v := range reserve {
			s.Repository.AddSkuStock(ctxTx, v.Sku, v.Amount, v.WarehouseID)
		}
		err = s.Repository.DeleteStocksReserveByOrderId(ctxTx, req.OrderId)
		if err != nil {
			return err
		}
		res, err = s.Repository.CancelOrder(ctxTx, req)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("cancel order: %w", err)
	}
	return res, nil
}

func (s *service) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	return s.Repository.Stocks(ctx, req)
}

// возвращает мапу: склад-количество, откуда брать нужное коль-во ску
func (s *service) ReserveStock(ctx context.Context, sku int64, count int64) (map[int64]int64, error) {
	var (
		warehouseIdReserveCnt = make(map[int64]int64, 1)
	)
	stocks, err := s.Repository.Stocks(ctx, &loms.StocksRequest{Sku: uint32(sku)})
	if err != nil {
		return nil, err
	}

	var (
		reservedCount int64
	)

	for _, v := range stocks.Stocks {
		warehouseID := v.WarehouseId
		warehouseStock := int64(v.Count)
		left := count - reservedCount
		if left == 0 {
			break
		}
		if warehouseStock >= left {
			warehouseIdReserveCnt[warehouseID] = left
			reservedCount += left
		} else {
			warehouseIdReserveCnt[warehouseID] = warehouseStock
			reservedCount += warehouseStock
		}
	}

	if reservedCount != count {
		return nil, errors.New("not enough stocks")
	}

	return warehouseIdReserveCnt, nil
}
