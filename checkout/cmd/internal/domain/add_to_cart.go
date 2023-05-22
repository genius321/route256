package domain

import (
	"errors"
	"fmt"
)

var (
	ErrStockInsufficient = errors.New("stock insufficient")
)

func (m *Model) AddToCart(_ int64, sku uint32, count uint16) error {
	stocks, err := m.stockChecker.Stocks(sku)
	if err != nil {
		return fmt.Errorf("get stocks: %w", err)
	}

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter == 0 {
			return nil
		}
	}

	return ErrStockInsufficient
}
