package loms

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/clientwrapper"
)

type CreateOrderRequest struct {
	User  int64             `json:"user"`
	Items []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type CreateOrderResponse struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) CreateOrder(ctx context.Context, user int64) (*domain.Order, error) {
	request := CreateOrderRequest{
		User: user,
		Items: []CreateOrderItem{
			{
				SKU:   1076963,
				Count: 1,
			},
		},
	}

	response, err := clientwrapper.Do[CreateOrderRequest, CreateOrderResponse](ctx, c.urlCreateOrder, request)
	if err != nil {
		return nil, err
	}

	return &domain.Order{
		OrderID: response.OrderID,
	}, nil
}
