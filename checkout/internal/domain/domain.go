package domain

import "context"

type StockChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Model struct {
	stockChecker StockChecker
}

func New(stockChecker StockChecker) *Model {
	return &Model{stockChecker: stockChecker}
}
