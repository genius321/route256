package purchase

import (
	"log"
	"route256/checkout/internal/domain"
)

type Handler struct {
	Model *domain.ModelOrderCreater
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Request struct {
	User int64 `json:"user"`
}

func (h Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	err := h.Model.Purchase(req.User)
	return Response{OrderID: 17}, err
}
