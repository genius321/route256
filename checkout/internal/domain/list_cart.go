package domain

import (
	"fmt"
	"log"
)

type ProductGetter interface {
	GetProduct(token string, sku int64) (Product, error)
}

type ModelProductGetter struct {
	productGetter ProductGetter
}

func NewModelProductGetter(productGetter ProductGetter) *ModelProductGetter {
	return &ModelProductGetter{productGetter: productGetter}
}

func (m *ModelProductGetter) GetProduct(token string, sku uint32) (Product, error) {
	product, err := m.productGetter.GetProduct(token, int64(sku))
	if err != nil {
		return product, fmt.Errorf("create order: %w", err)
	}
	log.Printf("product: %v", product)

	return product, nil
}
