package addtocart

import (
	"errors"
	"log"
	"route256/checkout/internal/domain"
)

type Handler struct {
	Model *domain.ModelStockChecker
}

type Response struct {
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (h Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	err := h.Model.AddToCart(req.User, req.SKU, req.Count)
	return Response{}, err
}
