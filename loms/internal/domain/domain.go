package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type LomsRepository interface {
	CreateOrder(context.Context, *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	ReserveItems(context.Context, *model.CreateOrderRequest) error
	ListOrder(context.Context, *model.ListOrderRequest) (*model.ListOrderResponse, error)
	//OrderPayed
	//CancelOrder
	Stocks(context.Context, *model.StocksRequest) (*model.StocksResponse, error)
}

type repository struct {
	lomsRepository     LomsRepository
	transactionManager TransactionManager
}

func NewRepository(lomsRepository LomsRepository, transactionManager TransactionManager) *repository {
	return &repository{
		lomsRepository:     lomsRepository,
		transactionManager: transactionManager,
	}
}

var _ Service = (*service)(nil)

type Service interface {
	CreateOrder(context.Context, *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	ListOrder(context.Context, *model.ListOrderRequest) (*model.ListOrderResponse, error)
	OrderPayed(context.Context, *model.OrderPayedRequest) error
	CancelOrder(context.Context, *model.CancelOrderRequest) error
	Stocks(context.Context, *model.StocksRequest) (*model.StocksResponse, error)
}

type service struct {
	repository *repository
}

func NewService(repository *repository) *service {
	return &service{
		repository: repository,
	}
}
