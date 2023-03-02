package domain

import (
	"context"
	loms "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
)

type AddToCartRequest struct {
	// user ID
	User int64
	// stock keeping unit - единица складского учёта
	Sku   uint32
	Count uint32
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (m *model) AddToCart(ctx context.Context, req *AddToCartRequest) error {
	stocks, err := m.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
	if err != nil {
		return errors.WithMessage(err, "checking stocks")
	}

	counter := int64(req.Count)
	for _, stock := range stocks.GetStocks() {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return nil
		}
	}

	return ErrInsufficientStocks
}
