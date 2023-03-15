package cancel_order

import (
	"route256/loms/internal/domain"
)

type CancelOrderDaemon interface {
}

type cancelOrderDaemon struct {
	businessLogic domain.Service
}

func NewCancelOrderDaemon(businessLogic domain.Service) *cancelOrderDaemon {
	return &cancelOrderDaemon{
		businessLogic: businessLogic,
	}
}

func runCancelDaemon(businessLogic domain.Service) {
	// создаём функцию на обработку

	// запускаем воркеров

	// по тикеру раз в минуту проверяем заказы на отмену

	// отменяем заказы (работа для воркеров)

}
