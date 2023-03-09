package domain

import (
	"context"
	"route256/checkout/internal/domain/model"
)

func (s *service) DeleteFromCart(ctx context.Context, req *model.DeleteFromCartRequest) error {
	// удаляем товар из корзины

	err := s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// проверяем, полностью надо всё удалить или только часть
		count, err := s.repository.checkoutRepository.GetCartItemCount(ctxTX, req.User, req.Sku)
		if err != nil {
			return err
		}

		part := (req.Count < count)

		err = s.repository.checkoutRepository.DeleteFromCart(ctxTX, part, req)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
