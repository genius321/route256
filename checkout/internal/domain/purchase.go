package domain

import (
	"fmt"
	"log"
)

type OrderCreater interface {
	CreateOrder(user int64, items []Item) (int64, error)
}

type ModelOrderCreater struct {
	orderCreater OrderCreater
}

func NewModelOrderCreater(orderCreater OrderCreater) *ModelOrderCreater {
	return &ModelOrderCreater{orderCreater: orderCreater}
}

func (m *ModelOrderCreater) Purchase(user int64) error {
	items := []Item{{SKU: 4487693, Count: 15}, {SKU: 32956725, Count: 31}}
	orderID, err := m.orderCreater.CreateOrder(user, items)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	log.Printf("orderID: %v", orderID)

	return nil
}
