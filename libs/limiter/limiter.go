package limiter

import (
	"context"
	"time"
)

type Limiter interface {
	Wait(context.Context) error
}

type limiter struct {
	// время освобождения одного запроса
	ticker *time.Ticker
	// канал для блокировки
	ch chan struct{}
}

func (l *limiter) run() {
	for {
		// по таймеру вычитываем из канала
		select {
		case <-l.ticker.C:
			// вычитываем из канала значение, освобождаем место под новый запрос
			<-l.ch
		}
	}
}

func (l *limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// можем записать в канал, значит есть свободные запросы, получаем управление
	case l.ch <- struct{}{}:
		return nil
	}
}

// создаёт лимитер
// d - частота запросов
// count - количество запросов
func NewLimiter(d time.Duration, count int) *limiter {
	l := &limiter{
		ticker: time.NewTicker(d / time.Duration(count)),
		ch:     make(chan struct{}, count),
	}
	go l.run()

	return l
}
