package stocks

import (
	"context"
	"log"
	"time"
)

type Handler struct {
}

type Response struct {
	Stocks []StockItem `json:"stocks"`
}

type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type Request struct {
	SKU uint32 `json:"sku"`
}

func (r Request) Validate() error {
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v", req)
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ctx.Done():
		log.Println(ctx.Err())
		return Response{}, ctx.Err()
	}
	return Response{
		Stocks: []StockItem{
			{WarehouseID: 1, Count: 200},
		},
	}, nil
}
