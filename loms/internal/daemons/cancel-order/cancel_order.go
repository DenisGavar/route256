package cancel_order

import (
	"context"
	"log"
	workerpool "route256/libs/worker-pool"
	"route256/loms/internal/domain/model"
	"time"

	"github.com/pkg/errors"
)

// должен мочь получить список заказов на отмену (работа с репозиторием, но через сервис)

// должен мочь отменить заказ по ID

type CancelOrderDaemon interface {
	RunCancelDaemon(workersCount int, cancelOrderTime time.Duration)
}

type OrderCanceler interface {
	OrdersToCancel(context.Context, time.Time) ([]*model.CancelOrderRequest, error)
	CancelOrder(context.Context, *model.CancelOrderRequest) error
}

type cancelOrderDaemon struct {
	orderCanceler OrderCanceler
}

func NewCancelOrderDaemon(orderCanceler OrderCanceler) *cancelOrderDaemon {
	return &cancelOrderDaemon{
		orderCanceler: orderCanceler,
	}
}

var _ CancelOrderDaemon = (*cancelOrderDaemon)(nil)

func (c *cancelOrderDaemon) RunCancelDaemon(workersCount int, cancelOrderTime time.Duration) {
	ctx, cacnel := context.WithCancel(context.Background())
	defer cacnel()

	ticker := time.NewTicker(time.Minute * 1)

	// создаём функцию на обработку
	callback := func(cancelOrderRequest *model.CancelOrderRequest) *workerpool.Result[*model.CancelOrderRequest] {
		log.Println("daemon: cancelling order")

		err := c.orderCanceler.CancelOrder(ctx, cancelOrderRequest)
		if err != nil {
			// если ошибка при получении данных, то возвращаем ошибку
			return &workerpool.Result[*model.CancelOrderRequest]{
				Out:   nil,
				Error: errors.WithMessage(err, "cancelling order"),
			}
		}

		// возвращаем ту же самую структуру, чтобы было понятно, какой заказ отменили
		return &workerpool.Result[*model.CancelOrderRequest]{
			Out:   cancelOrderRequest,
			Error: nil,
		}
	}

	// запускаем воркеров
	pool, results := workerpool.New[*model.CancelOrderRequest, *model.CancelOrderRequest](ctx, workersCount)
	pool.Init(ctx)

	go func() {
		// читаем из канала пока он открыт
		for result := range results {
			if result.Error != nil {
				// ошибку логируем
				log.Println(result.Error)
			}
			// отмечаем, что задача выполнена, результат получен
			pool.JobDone()
		}
	}()

	// по тикеру раз в минуту проверяем заказы на отмену

	for {
		select {
		case <-ticker.C:
			// получаем заказы на отмену, передаём время, с которого заказы надо отменять
			log.Println("checking orders to cancel")
			ordersToCancel, err := c.orderCanceler.OrdersToCancel(ctx, time.Now().Add(-cancelOrderTime))
			if err != nil {
				// ошибку логируем
				log.Println(err)
				break
			}
			for _, orderToCancel := range ordersToCancel {
				// отменяем заказы (работа для воркеров)
				pool.Push(ctx, &workerpool.Job[*model.CancelOrderRequest, *model.CancelOrderRequest]{
					Callback: callback,
					Args:     orderToCancel,
				})
			}
		case <-ctx.Done():
			// вышли по отмене контекста
			// останавливаем worker pool, передаём аргумент, что ждём завершения всех отправленных задач
			pool.Stop(true)
			return
		}
	}

}
