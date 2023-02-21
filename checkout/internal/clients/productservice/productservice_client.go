package productservice

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/clientwrapper"
)

type Client struct {
	url   string
	token string

	urlGetProduct string
}

func New(url string, token string) *Client {
	return &Client{
		url:   url,
		token: token,

		urlGetProduct: url + "/get_product",
	}
}

type GetProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, sku uint32) (*domain.Product, error) {
	request := GetProductRequest{
		Token: c.token,
		SKU:   sku,
	}

	response, err := clientwrapper.Do[GetProductRequest, GetProductResponse](ctx, c.urlGetProduct, request)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}
