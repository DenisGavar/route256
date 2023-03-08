package repository

import (
	"route256/libs/transactor"
)

const (
	itemsStocksTable            = "items_stocks"
	ordersTable                 = "orders"
	orderItemsTable             = "order_items"
	itemsStocksReservationTable = "items_stocks_reservation"
)

type repo struct {
	queryEngineProvider transactor.QueryEngineProvider
}

func NewRepo(queryEngineProvider transactor.QueryEngineProvider) *repo {
	return &repo{
		queryEngineProvider: queryEngineProvider,
	}
}
