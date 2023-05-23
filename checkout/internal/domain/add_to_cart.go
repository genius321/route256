package domain

import (
	"errors"
	"fmt"
	"log"
)

type StockChecker interface {
	Stocks(sku uint32) ([]Stock, error)
}

type ModelStockChecker struct {
	stockChecker StockChecker
}

func NewModelStockChecker(stockChecker StockChecker) *ModelStockChecker {
	return &ModelStockChecker{stockChecker: stockChecker}
}

var (
	ErrStockInsufficient = errors.New("stock insufficient")
)

func (m *ModelStockChecker) AddToCart(_ int64, sku uint32, count uint16) error {
	stocks, err := m.stockChecker.Stocks(sku)
	if err != nil {
		return fmt.Errorf("get stocks: %w", err)
	}

	log.Printf("stocks: %v", stocks)

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return nil
		}
	}

	return ErrStockInsufficient
}
