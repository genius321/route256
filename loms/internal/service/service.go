package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	orderModels "route256/loms/internal/models/order"
	stockModels "route256/loms/internal/models/stock"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Repository interface {
	CreateOrder(context.Context, orderModels.User, orderModels.Items) (orderModels.OrderId, error)
	ListOrder(context.Context, orderModels.OrderId) (orderModels.Status, orderModels.User, orderModels.Items, error)

	Stocks(context.Context, stockModels.Sku) (stockModels.Stocks, error)
	TakeSkuStock(ctx context.Context, sku int64, amount int64, warehouseID int64) (int64, error)
	AddSkuStockReserve(ctx context.Context, sku int64, amount int64, warehouseID int64, orderId int64) error
	// AddSkuStock(ctx context.Context, sku int64, amount int64, warehouseID int64) error
}

// ничего не знает про транспортный уровень,
// но знает, что бд и тракнзакционный мендеджер реализуют необходимое поведение
type Service struct {
	Repository
	TransactionManager
}

func NewService(r Repository, tm TransactionManager) *Service {
	return &Service{Repository: r, TransactionManager: tm}
}

func (s *Service) CreateOrder(
	ctx context.Context,
	user orderModels.User,
	items orderModels.Items,
) (orderModels.OrderId, error) {
	var orderId orderModels.OrderId
	var err error
	err = s.RunSerializable(ctx, func(ctxTx context.Context) error {
		orderId, err = s.Repository.CreateOrder(ctxTx, user, items)
		if err != nil {
			return err
		}
		for _, v := range items {
			warehouseIdReserveCnt, err := s.ReserveStock(ctxTx, stockModels.Sku(v.Sku), stockModels.Count(v.Count))
			if err != nil {
				return err
			}
			log.Println(warehouseIdReserveCnt)
			for warehouseID, reserveCnt := range warehouseIdReserveCnt {
				_, err = s.Repository.TakeSkuStock(ctxTx, int64(v.Sku), reserveCnt, warehouseID)
				if err != nil {
					return err
				}
				err = s.Repository.AddSkuStockReserve(ctxTx, int64(v.Sku), reserveCnt, warehouseID, int64(orderId))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("create order: %w", err)
	}
	log.Println(orderId)
	return orderId, nil
}

type mWarehouseIdReserveCnt map[int64]int64

func (s *Service) ReserveStock(
	ctx context.Context,
	sku stockModels.Sku,
	count stockModels.Count,
) (mWarehouseIdReserveCnt, error) {
	var warehouseIdReserveCnt = make(mWarehouseIdReserveCnt, 1)

	stocks, err := s.Repository.Stocks(ctx, sku)
	if err != nil {
		return nil, err
	}

	var reservedCount int64

	for _, v := range stocks {
		warehouseID := int64(v.WarehouseId)
		warehouseStock := int64(v.Count)
		if warehouseStock == 0 {
			continue
		}

		left := int64(count) - reservedCount
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

	if reservedCount != int64(count) {
		return nil, errors.New("not enough stocks")
	}

	return warehouseIdReserveCnt, nil
}

func (s *Service) Stocks(ctx context.Context, sku stockModels.Sku) (stockModels.Stocks, error) {
	return s.Repository.Stocks(ctx, sku)
}

func (s *Service) ListOrder(
	ctx context.Context,
	orderId orderModels.OrderId,
) (orderModels.Status, orderModels.User, orderModels.Items, error) {
	return s.Repository.ListOrder(ctx, orderId)
}
