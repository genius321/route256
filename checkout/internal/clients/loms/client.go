package loms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"route256/checkout/internal/domain"
	"time"
)

const (
	StocksPath = "stocks"
)

type Client struct {
	pathStock string
}

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksResponse struct {
	Stocks []struct {
		WarehouseID int64  `json:"warehouseID"`
		Count       uint64 `json:"count"`
	} `json:"stocks"`
}

func New(clientUrl string) *Client {
	stockUrl, _ := url.JoinPath(clientUrl, StocksPath)
	return &Client{pathStock: stockUrl}
	// return &Client{pathStock: fmt.Sprintf("%s/%s", clientUrl, StocksPath)}
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	requestStocks := StocksRequest{SKU: sku}

	rawData, err := json.Marshal(&requestStocks)
	if err != nil {
		return nil, fmt.Errorf("encode stock request: %w", err)
	}

	ctx, fnCancel := context.WithTimeout(ctx, 5*time.Second)
	defer fnCancel()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.pathStock, bytes.NewBuffer(rawData))
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
