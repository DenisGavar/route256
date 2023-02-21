package domain

import "context"

type StocksChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type ProductGetter interface {
	GetProduct(ctx context.Context, sku uint32) (*Product, error)
}

type Model struct {
	stocksChecker StocksChecker
	productGetter ProductGetter
}

func New(stocksChecker StocksChecker, productGetter ProductGetter) *Model {
	return &Model{
		stocksChecker: stocksChecker,
		productGetter: productGetter,
	}
}
