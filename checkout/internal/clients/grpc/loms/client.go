package loms

import (
	"context"
	lomsServiceAPI "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
)

type LomsClient interface {
	Stocks(context.Context, *lomsServiceAPI.StocksRequest) (*lomsServiceAPI.StocksResponse, error)
	CreateOrder(context.Context, *lomsServiceAPI.CreateOrderRequest) (*lomsServiceAPI.CreateOrderResponse, error)
}

type client struct {
	lomsClient lomsServiceAPI.LOMSV1Client
}

func New(cc *grpc.ClientConn) *client {
	return &client{
		lomsClient: lomsServiceAPI.NewLOMSV1Client(cc),
	}
}
