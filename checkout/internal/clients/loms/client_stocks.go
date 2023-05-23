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
	StocksPath = "stocks"
)

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksResponse struct {
	Stocks []struct {
		WarehouseID int64  `json:"warehouseID"`
		Count       uint64 `json:"count"`
	} `json:"stocks"`
}

type ClientStocks struct {
	pathStock string
}

func NewClientStocks(clientUrl string) *ClientStocks {
	stockUrl, _ := url.JoinPath(clientUrl, StocksPath)
	return &ClientStocks{pathStock: stockUrl}
}

func (c *ClientStocks) Stocks(sku uint32) ([]domain.Stock, error) {
	requestStocks := StocksRequest{SKU: sku}

	rawData, err := json.Marshal(&requestStocks)
	if err != nil {
		return nil, fmt.Errorf("encode stock request: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, c.pathStock, bytes.NewBuffer(rawData))
	if err != nil {
		return nil, fmt.Errorf("prepare stock request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("do stock request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code get stock: %d", httpResponse.StatusCode)
	}

	responseStocks := StocksResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responseStocks)
	if err != nil {
		return nil, fmt.Errorf("decode stock request: %w", err)
	}

	result := make([]domain.Stock, 0, len(responseStocks.Stocks))
	for _, v := range responseStocks.Stocks {
		result = append(result, domain.Stock{
			WarehouseID: v.WarehouseID,
			Count:       v.Count,
		})
	}

	return result, nil
}
