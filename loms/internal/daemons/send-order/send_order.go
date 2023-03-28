package send_order

import (
	"context"
	"log"
	workerPool "route256/libs/worker-pool"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/sender"
	"time"

	"github.com/pkg/errors"
)

type SendOrderDaemon interface {
	RunSendDaemon(workersCount int, topic string)
}

type Messager interface {
	// должен мочь получить из репозитория то, что надо отправлять
	MessagesToSend(context.Context) ([]*model.OrderMessage, error)
	// должен мочь пометить отправленное сообщение
	MessageSent(context.Context, int64) error
}

type OrderSender interface {
	SendOrder(context.Context, *sender.OrderMessage) error
}

type sendOrderDaemon struct {
	messager    Messager
	orderSender OrderSender
}

func NewSendOrderDaemon(messager Messager, orderSender OrderSender) *sendOrderDaemon {
	return &sendOrderDaemon{
		messager:    messager,
		orderSender: orderSender,
	}
}

var _ SendOrderDaemon = (*sendOrderDaemon)(nil)

func (c *sendOrderDaemon) RunSendDaemon(workersCount int, topic string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// по тикеру запускаем обработку отправки
	ticker := time.NewTicker(time.Second * 30)

	// создаём функцию на обработку
	// возвращаем id записи в БД
	callback := func(orderMessage *sender.OrderMessage) *workerPool.Result[int64] {
		log.Println("daemon: sending order")

		// отправялем сообщение
		err := c.orderSender.SendOrder(ctx, orderMessage)
		if err != nil {
			// если ошибка при получении данных, то возвращаем ошибку
			return &workerPool.Result[int64]{
				Out:   orderMessage.OutboxKey,
				Error: errors.WithMessage(err, "sending order"),
			}
		}

		// возвращаем ключ. чтобы было понятно, что отправили
		return &workerPool.Result[int64]{
			Out:   orderMessage.OutboxKey,
			Error: nil,
		}
	}

	// запускаем воркеров
	pool := workerPool.New[*sender.OrderMessage, int64](ctx, workersCount)
	pool.Init(ctx)

	// создаём канал для чтения результатов, буфер на количество воркеров
	results := make(chan *workerPool.Result[int64], workersCount)
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
			// помечаем в БД сообщение отправленным
			err := c.messager.MessageSent(ctx, result.Out)
			if err != nil {
				log.Println(err)
			}
			// отмечаем, что задача выполнена, результат получен
			pool.JobDone()
		}
	}()

	// по тикеру проверяем сообщения для отправки

	for {
		select {
		case <-ticker.C:
			// получаем сообщения для отправки
			log.Println("checking messages to send")
			messagesToSend, err := c.messager.MessagesToSend(ctx)
			if err != nil {
				// ошибку логируем
				log.Println(err)
				break
			}
			for _, messageToSend := range messagesToSend {
				// отправляем сообщения (работа для воркеров)
				pool.Push(ctx, &workerPool.Job[*sender.OrderMessage, int64]{
					Callback: callback,
					Args: &sender.OrderMessage{
						OutboxKey: messageToSend.Id,
						Key:       messageToSend.OrderId,
						Message:   messageToSend.Message,
						Timestamp: messageToSend.CreatedAt,
						Topic:     topic,
					},
					Results: results,
				})
			}
		case <-ctx.Done():
			// вышли по отмене контекста
			return
		}
	}

}
