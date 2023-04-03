package domain

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	"route256/libs/logger"
	loms "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) AddToCart(ctx context.Context, req *model.AddToCartRequest) error {
	// добавляем товар в корзину
	logger.Debug("checkout domain", zap.String("handler", "AddToCart"), zap.String("request", fmt.Sprintf("%+v", req)))

	// проверяем, что товара достаточно на складах
	stocks, err := s.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
	if err != nil {
		return errors.WithMessage(err, ErrCheckingStocks.Error())
	}

	counter := int64(req.Count)
	for _, stocksItem := range stocks.GetStocks() {
		counter -= int64(stocksItem.Count)
		if counter <= 0 {
			// если товаров достаточно, что добавляем в корзину
			err = s.repository.checkoutRepository.AddToCart(ctx, req)
			if err != nil {
				return errors.WithMessage(err, ErrAddingToCart.Error())
			}
			return nil
		}
	}

	return ErrNotEnoughItems
}
