package createorder

import (
	"log"
)

type Handler struct {
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Request struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (h *Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	return Response{OrderID: 666}, nil
}
