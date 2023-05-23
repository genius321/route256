package orderpayed

import (
	"log"
)

type Handler struct {
}

type Response struct {
}

type Request struct {
	OrderID int64 `json:"orderID"`
}

func (h *Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	return Response{}, nil
}
