package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	orderModels "route256/loms/internal/models/order"
	stockModels "route256/loms/internal/models/stock"
)

type TransactionManager interface {
	RunRepeatableRead(context.Context, func(ctxTx context.Context) error) error
	RunSerializable(context.Context, func(ctxTx context.Context) error) error
}

type Repository interface {
	CreateOrder(context.Context, orderModels.User, orderModels.Items) (orderModels.OrderId, error)
	ListOrder(context.Context, orderModels.OrderId) (
		orderModels.Status, orderModels.User, orderModels.Items, error)
	OrderPayed(context.Context, orderModels.OrderId) error
	CancelOrder(context.Context, orderModels.OrderId) error

	Stocks(context.Context, stockModels.Sku) (stockModels.Stocks, error)
	TakeSkuStock(context.Context, stockModels.StockWithSku) (stockModels.Count, error)
	AddSkuStockReserve(context.Context, stockModels.StockWithSku, orderModels.OrderId) error
	DeleteStocksReserveByOrderId(context.Context, orderModels.OrderId) error
	TakeStocksReserveByOrderId(context.Context, orderModels.OrderId) (stockModels.StocksWithSku, error)
	AddSkuStock(context.Context, stockModels.StockWithSku) error
}

// ничего не знает про транспортный уровень,
// но знает, что бд и тракнзакционный мендеджер реализуют необходимое поведение
type Business struct {
	Repository
	TransactionManager
}

func NewBusiness(r Repository, tm TransactionManager) *Business {
	return &Business{Repository: r, TransactionManager: tm}
}

func (s *Business) CreateOrder(
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
				_, err = s.Repository.TakeSkuStock(ctxTx, stockModels.StockWithSku{
					Sku: stockModels.Sku(v.Sku),
					Stock: stockModels.Stock{
						WarehouseId: stockModels.WarehouseId(warehouseID),
						Count:       stockModels.Count(reserveCnt),
					},
				})
				if err != nil {
					return err
				}
				err = s.Repository.AddSkuStockReserve(
					ctxTx,
					stockModels.StockWithSku{
						Sku: stockModels.Sku(v.Sku),
						Stock: stockModels.Stock{
							WarehouseId: stockModels.WarehouseId(warehouseID),
							Count:       stockModels.Count(reserveCnt),
						},
					},
					orderId)
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

func (s *Business) ReserveStock(
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

func (s *Business) Stocks(ctx context.Context, sku stockModels.Sku) (stockModels.Stocks, error) {
	return s.Repository.Stocks(ctx, sku)
}

func (s *Business) ListOrder(
	ctx context.Context,
	orderId orderModels.OrderId,
) (orderModels.Status, orderModels.User, orderModels.Items, error) {
	return s.Repository.ListOrder(ctx, orderId)
}

func (s *Business) OrderPayed(ctx context.Context, orderId orderModels.OrderId) error {
	err := s.RunSerializable(ctx, func(ctxTx context.Context) error {
		err := s.Repository.DeleteStocksReserveByOrderId(ctxTx, orderId)
		if err != nil {
			return err
		}
		err = s.Repository.OrderPayed(ctxTx, orderId)
		return err
	})
	if err != nil {
		return fmt.Errorf("order payed: %w", err)
	}
	return nil
}

func (s *Business) CancelOrder(ctx context.Context, orderId orderModels.OrderId) error {
	err := s.RunSerializable(ctx, func(ctxTx context.Context) error {
		reserve, err := s.Repository.TakeStocksReserveByOrderId(ctxTx, orderId)
		if err != nil {
			return err
		}
		for _, v := range reserve {
			s.Repository.AddSkuStock(ctxTx, v)
		}
		err = s.Repository.DeleteStocksReserveByOrderId(ctxTx, orderId)
		if err != nil {
			return err
		}
		err = s.Repository.CancelOrder(ctxTx, orderId)
		return err
	})
	if err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}
	return nil
}
