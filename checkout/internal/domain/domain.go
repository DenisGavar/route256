package domain

import "context"

type StocksChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type ProductGetter interface {
	GetProduct(ctx context.Context, sku uint32) (*Product, error)
}

type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64) (*Order, error)
}

type Model struct {
	stocksChecker StocksChecker
	productGetter ProductGetter
	orderCreator  OrderCreator
}

func New(stocksChecker StocksChecker, orderCreator OrderCreator, productGetter ProductGetter) *Model {
	return &Model{
		stocksChecker: stocksChecker,
		orderCreator:  orderCreator,
		productGetter: productGetter,
	}
}
