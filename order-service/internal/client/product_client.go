package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/dto"
)

type ProductClient struct {
	baseURL string
	http    *http.Client
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{baseURL: baseURL, http: &http.Client{}}
}

func (client *ProductClient) GetProducts(ctx context.Context, jwtCookie *http.Cookie) ([]dto.ProductResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseURL+"/products", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(jwtCookie)
	resp, err := client.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product-service responde status %d", resp.StatusCode)
	}

	var products []dto.ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
