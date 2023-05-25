package productservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"route256/checkout/internal/domain"
)

const (
	GetProductPath = "get_product"
)

type ClientGetProduct struct {
	pathGetProduct string
}

type GetProductRequest struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

func NewClientGetProduct(clientUrl string) *ClientGetProduct {
	getProductUrl, _ := url.JoinPath(clientUrl, GetProductPath)
	return &ClientGetProduct{pathGetProduct: getProductUrl}
}

func (c *ClientGetProduct) GetProduct(token string, sku int64) (domain.Product, error) {
	var result domain.Product

	rawData, err := json.Marshal(&GetProductRequest{Token: token, SKU: sku})
	if err != nil {
		return result, fmt.Errorf("encode getproduct request: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, c.pathGetProduct, bytes.NewBuffer(rawData))
	if err != nil {
		return result, fmt.Errorf("prepare getproduct request: %w", err)
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return result, fmt.Errorf("do getproduct request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return result, fmt.Errorf("wrong status code getproduct: %d", httpResponse.StatusCode)
	}

	responseGetProduct := GetProductResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responseGetProduct)
	if err != nil {
		return result, fmt.Errorf("decode getproduct request: %w", err)
	}

	return domain.Product{
		Name:  responseGetProduct.Name,
		Price: uint32(responseGetProduct.Price),
	}, nil
}
