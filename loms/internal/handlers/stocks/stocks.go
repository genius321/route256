package stocks

import (
	"log"
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

func (h *Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	return Response{
		Stocks: []StockItem{
			{WarehouseID: 1, Count: 150},
			{WarehouseID: 2, Count: 50},
		},
	}, nil
}
