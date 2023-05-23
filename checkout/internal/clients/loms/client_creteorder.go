package loms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"route256/checkout/internal/domain"
)

const (
	CreateOrderPath = "createOrder"
)

type CreateOrderRequest struct {
	User  int64         `json:"user"`
	Items []domain.Item `json:"items"`
}

type CreateOrderResponse struct {
	OrderID int64 `json:"orderID"`
}

type ClientCreateOrder struct {
	pathCreateOrder string
}

func NewClientCreateOrder(clientUrl string) *ClientCreateOrder {
	createOrderUrl, _ := url.JoinPath(clientUrl, CreateOrderPath)
	return &ClientCreateOrder{pathCreateOrder: createOrderUrl}
}

func (c *ClientCreateOrder) CreateOrder(user int64, items []domain.Item) (int64, error) {
	requestCreateOrder := CreateOrderRequest{User: user, Items: items}

	rawData, err := json.Marshal(&requestCreateOrder)
	if err != nil {
		return 0, fmt.Errorf("encode createorder request: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, c.pathCreateOrder, bytes.NewBuffer(rawData))
	if err != nil {
		return 0, fmt.Errorf("prepare createorder request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return 0, fmt.Errorf("do createorder request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("wrong status code createorder: %d", httpResponse.StatusCode)
	}

	responseCreateOrder := CreateOrderResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responseCreateOrder)
	if err != nil {
		return 0, fmt.Errorf("decode createorder request: %w", err)
	}

	result := responseCreateOrder.OrderID

	return result, nil
}
