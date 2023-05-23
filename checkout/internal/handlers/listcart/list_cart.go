package listcart

import (
	"errors"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
)

type Handler struct {
	Model *domain.ModelProductGetter
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
	items := []domain.Item{{SKU: 4487693, Count: 15}, {SKU: 32956725, Count: 31}}
	resp := []Item{}
	var totalprice uint32
	for _, v := range items {
		product, err := h.Model.GetProduct(config.AppConfig.Token, v.SKU)
		if err != nil {
			return Response{}, err
		}
		resp = append(resp, Item{SKU: v.SKU, Count: v.Count, Name: product.Name, Price: product.Price})
		totalprice += product.Price * uint32(v.Count)
	}
	return Response{Items: resp,
		TotalPrice: totalprice,
	}, nil
}
