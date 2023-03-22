package repository

import (
	"context"
	"route256/libs/transactor"
	"route256/loms/internal/domain/model"
	"time"
)

const (
	itemsStocksTable            = "items_stocks"
	ordersTable                 = "orders"
	orderItemsTable             = "order_items"
	itemsStocksReservationTable = "items_stocks_reservation"
)

//go:generate mockgen -source="repository.go" -destination="mocks/repository_mock.go" -package=mocks . LomsRepository
type LomsRepository interface {
	CreateOrder(ctx context.Context, createOrderRequest *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	ListOrder(ctx context.Context, listOrderRequest *model.ListOrderRequest) (*model.ListOrderResponse, error)
	ClearReserves(ctx context.Context, orderId int64) error
	Reserves(ctx context.Context, orderId int64) (*model.Reserve, error)
	ReturnReserve(ctx context.Context, reserveStocksItem *model.ReserveStocksItem) error
	Stocks(ctx context.Context, stocksRequest *model.StocksRequest) (*model.StocksResponse, error)
	ReserveItems(ctx context.Context, orderId int64, req *model.ReserveStocksItem) error
	ChangeStatus(ctx context.Context, orderId int64, status string) error
	OrdersToCancel(ctx context.Context, time time.Time) ([]*model.CancelOrderRequest, error)
}

type repository struct {
	queryEngineProvider transactor.QueryEngineProvider
}

func NewRepository(queryEngineProvider transactor.QueryEngineProvider) *repository {
	return &repository{
		queryEngineProvider: queryEngineProvider,
	}
}
