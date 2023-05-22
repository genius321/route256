package listcart

import (
	"errors"
	"log"
)

type Handler struct {
}

type Response struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"totalprice"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Request struct {
	User int64 `json:"user"`
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
	return Response{}, nil
}
