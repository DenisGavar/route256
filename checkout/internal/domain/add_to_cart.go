package domain

import (
	"context"

	"github.com/pkg/errors"
)

type AddToCartRequest struct {
	// user ID
	User int64
	// stock keeping unit - единица складского учёта
	Sku   uint32
	Count uint32
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (m *model) AddToCart(ctx context.Context, req *AddToCartRequest) error {
	return nil

	// stocks, err := m.stocksChecker.Stocks(ctx, sku)
	// if err != nil {
	// 	return errors.WithMessage(err, "checking stocks")
	// }

	// counter := int64(count)
	// for _, stock := range stocks {
	// 	counter -= int64(stock.Count)
	// 	if counter <= 0 {
	// 		return nil
	// 	}
	// }

	// return ErrInsufficientStocks
}
