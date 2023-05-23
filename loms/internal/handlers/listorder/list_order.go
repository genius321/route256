package listorder

import (
	"log"
)

type Handler struct {
}

type Response struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Item `json:"items"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Request struct {
	OrderID int64 `json:"orderID"`
}

func (h *Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	return Response{
		Status: "new",
		User:   111,
		Items: []Item{
			{SKU: 222, Count: 333},
			{SKU: 444, Count: 555},
		},
	}, nil
}
