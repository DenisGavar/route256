package domain

import (
	"context"
	"route256/checkout/internal/domain/model"
	loms "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
)

func (s *service) AddToCart(ctx context.Context, req *model.AddToCartRequest) error {
	// добавляем товар в корзину

	// нулевое количество добавлять нет смысла
	if req.Count == 0 {
		return ErrNullCount
	}

	// проверяем, что товара достаточно на складах
	stocks, err := s.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
	if err != nil {
		return errors.WithMessage(err, "checking stocks")
	}

	counter := int64(req.Count)
	for _, stocksItem := range stocks.GetStocks() {
		counter -= int64(stocksItem.Count)
		if counter <= 0 {
			// если товаров достаточно, что добавляем в корзину
			err = s.repository.checkoutRepository.AddToCart(ctx, req)
			if err != nil {
				return errors.WithMessage(err, "adding to cart")
			}
			return nil
		}
	}

	return ErrNotEnoughItems
}
