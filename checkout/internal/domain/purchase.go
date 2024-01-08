package domain

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	"route256/libs/logger"
	loms "route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) Purchase(ctx context.Context, req *model.PurchaseRequest) (*model.PurchaseResponse, error) {
	// создаём заказ
	logger.Debug("checkout domain", zap.String("handler", "Purchase"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout domain Purchase processing")
	defer span.Finish()

	span.SetTag("user", req.User)

	var order *loms.CreateOrderResponse
	var err error

	createOrderRequest := &loms.CreateOrderRequest{
		User: req.User,
	}

	err = s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// получаем товары из корзины
		listCart, err := s.repository.checkoutRepository.ListCart(ctxTX, &model.ListCartRequest{User: req.User})
		if err != nil {
			return errors.WithMessage(err, ErrGettingListCart.Error())
		}

		// если корзина пустая, то выходим
		// возвращаем order_id = 0
		if len(listCart.Items) == 0 {
			return nil
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
			return errors.WithMessage(err, ErrCreatingOrder.Error())
		}

		// если заказ создали, то чистим корзину
		for _, cartItem := range listCart.Items {
			err = s.repository.checkoutRepository.DeleteFromCart(ctxTX, &model.DeleteFromCartRequest{
				User:  req.User,
				Sku:   cartItem.Sku,
				Count: cartItem.Count,
			})
			if err != nil {
				return errors.WithMessage(err, ErrDeletingFromCart.Error())
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.PurchaseResponse{OrderId: order.GetOrderId()}, nil
}
