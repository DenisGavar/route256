package domain

import (
	"context"
	"errors"
	"route256/libs/transactor"
	"route256/loms/internal/domain/model"
	repository "route256/loms/internal/repository/postgres"
)

var (
	ErrNotEnoughItems        = errors.New("not enough items")
	ErrCheckingReserves      = errors.New("checking reserves")
	ErrReturningReserves     = errors.New("returning reserves")
	ErrClearingReserves      = errors.New("clearing reserves")
	ErrChangingStatus        = errors.New("changing status")
	ErrCreatingOrder         = errors.New("creating order")
	ErrGettingStocks         = errors.New("getting stocks")
	ErrReservingItems        = errors.New("reserving items")
	ErrGettingListOrder      = errors.New("getting list order")
	ErrGettingOrdersToCancel = errors.New("getting orders to cancel")
	ErrGettingMessagesToSend = errors.New("getting messages to send")
	ErrChangingMessageSent   = errors.New("changing message sent")
)

type repo struct {
	lomsRepository     repository.LomsRepository
	transactionManager transactor.TransactionManager
}

func NewRepository(lomsRepository repository.LomsRepository, transactionManager transactor.TransactionManager) *repo {
	return &repo{
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
	repository *repo
}

func NewService(repo *repo) *service {
	return &service{
		repository: repo,
	}
}
