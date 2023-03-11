package domain

import (
	"context"
	"route256/checkout/internal/domain/model"

	"github.com/pkg/errors"
)

func (s *service) DeleteFromCart(ctx context.Context, req *model.DeleteFromCartRequest) error {
	// удаляем товар из корзины

	// нулевое количество удалять нет смысла
	if req.Count == 0 {
		return ErrNullCount
	}

	err := s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// получаем количество товара в корзине
		count, err := s.repository.checkoutRepository.GetCartItemCount(ctxTX, req.User, req.Sku)
		if err != nil {
			return errors.WithMessage(err, "getting cart item")
		}

		// проверяем, полностью надо всё удалить или только часть
		part := (req.Count < count)

		// непосредственно удаление товара в нужном количестве из корзины
		err = s.repository.checkoutRepository.DeleteFromCart(ctxTX, part, req)
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
