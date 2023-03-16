package domain

import (
	"context"
	"route256/checkout/internal/domain/model"

	"github.com/pkg/errors"
)

func (s *service) DeleteFromCart(ctx context.Context, req *model.DeleteFromCartRequest) error {
	// удаляем товар из корзины
	err := s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// удаление товара в нужном количестве из корзины
		err := s.repository.checkoutRepository.DeleteFromCart(ctxTX, req)
		if err != nil {
			return errors.WithMessage(err, "deleting from cart")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
