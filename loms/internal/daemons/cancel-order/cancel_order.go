package cancel_order

import (
	"context"
	"log"
	workerPool "route256/libs/worker-pool"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(time.Minute * 1)

	// создаём функцию на обработку
	callback := func(cancelOrderRequest *model.CancelOrderRequest) *workerPool.Result[*model.CancelOrderRequest] {
		log.Println("daemon: cancelling order")

		err := c.orderCanceler.CancelOrder(ctx, cancelOrderRequest)
		if err != nil {
			// если ошибка при получении данных, то возвращаем ошибку
			return &workerPool.Result[*model.CancelOrderRequest]{
				Out:   nil,
				Error: errors.WithMessage(err, "cancelling order"),
			}
		}

		// возвращаем ту же самую структуру, чтобы было понятно, какой заказ отменили
		return &workerPool.Result[*model.CancelOrderRequest]{
			Out:   cancelOrderRequest,
			Error: nil,
		}
	}

	// запускаем воркеров
	pool := workerPool.New[*model.CancelOrderRequest, *model.CancelOrderRequest](ctx, workersCount)
	pool.Init(ctx)

	// создаём канал для чтения результатов, буфер на количество воркеров
	results := make(chan *workerPool.Result[*model.CancelOrderRequest], workersCount)
	// закрываем канал результатов
	defer close(results)
	// останавливаем worker pool (закрываем канал задач)
	defer pool.Stop(false)

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
				pool.Push(ctx, &workerPool.Job[*model.CancelOrderRequest, *model.CancelOrderRequest]{
					Callback: callback,
					Args:     orderToCancel,
					Results:  results,
				})
			}
		case <-ctx.Done():
			// вышли по отмене контекста
			return
		}
	}

}
