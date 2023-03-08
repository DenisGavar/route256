package domain

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
)

var (
	ErrNotEnoughItems = errors.New("not enough items")
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type LomsRepository interface {
	CreateOrder(context.Context, *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	ListOrder(context.Context, *model.ListOrderRequest) (*model.ListOrderResponse, error)
	ClearReserves(ctx context.Context, orderId int64) error
	Reserves(ctx context.Context, orderId int64) (*model.Reserve, error)
	ReturnReserve(ctx context.Context, reserveStocksItem *model.ReserveStocksItem) error
	Stocks(context.Context, *model.StocksRequest) (*model.StocksResponse, error)

	ReserveItems(ctx context.Context, orderId int64, warehouseId int64, req *model.ReserveStocksItem) error
	ChangeStatus(ctx context.Context, orderId int64, status string) error
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
