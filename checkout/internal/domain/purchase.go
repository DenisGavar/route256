package domain

import (
	"context"
	"route256/checkout/internal/domain/model"
	loms "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
)

func (s *service) Purchase(ctx context.Context, req *model.PurchaseRequest) (*model.PurchaseResponse, error) {

	var order *loms.CreateOrderResponse
	var err error

	createOrderRequest := &loms.CreateOrderRequest{
		User: req.User,
	}

	err = s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// получаем товары из корзины
		listCart, err := s.repository.checkoutRepository.ListCart(ctx, &model.ListCartRequest{User: req.User})
		if err != nil {
			return err
		}

		// дополняем структуру товарами
		items := make([]*loms.OrderItem, 0)
		for _, cartItem := range listCart.Items {
			item := &loms.OrderItem{
				Sku:   cartItem.Sku,
				Count: cartItem.Count,
			}
			items = append(items, item)
		}
		createOrderRequest.Items = items

		// создаём заказ
		order, err = s.lomsClient.CreateOrder(ctxTX, createOrderRequest)
		if err != nil {
			return errors.WithMessage(err, "create order")
		}

		// если заказ создали, то чистим корзину
		for _, cartItem := range listCart.Items {
			err = s.repository.checkoutRepository.DeleteFromCart(ctxTX, false, &model.DeleteFromCartRequest{
				User:  req.User,
				Sku:   cartItem.Sku,
				Count: cartItem.Count,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.PurchaseResponse{OrderId: order.GetOrderId()}, nil
}
