package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

var _ Service = (*service)(nil)

type Service interface {
	CreateOrder(context.Context, *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	ListOrder(context.Context, *model.ListOrderRequest) (*model.ListOrderResponse, error)
	OrderPayed(context.Context, *model.OrderPayedRequest) error
	CancelOrder(context.Context, *model.CancelOrderRequest) error
	Stocks(context.Context, *model.StocksRequest) (*model.StocksResponse, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}
