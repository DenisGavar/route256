package limiter

import (
	"context"
	"time"
)

//go:generate mockgen -source="limiter.go" -destination="mocks/limiter_mock.go" -package=mocks . Limiter
type Limiter interface {
	Wait(context.Context) error
}

type limiter struct {
	// допустимый "всплеск" запросов в единицу времени
	maxCount int
	// доступных запросов в еденицу времени
	count int
	// время освобождения одного запроса
	ticker *time.Ticker
	// канал для блокировки
	ch chan struct{}
}

func (l *limiter) run() {
	for {
		// если свободных запросов не осталось, ждём пока не появятся
		if l.count <= 0 {
			// ждём ответа от тикера
			<-l.ticker.C
			// увеличиваем счетчик свободных запросов, но не больше максимума
			l.count = min(l.maxCount, l.count+1)
		}

		// если есть свободные запросы
		select {
		// либо пишем в канал и возвращаем управление
		case l.ch <- struct{}{}:
			// уменьшаем количество свободных запросов
			l.count--
		// либо ждём ответа от тикера
		case <-l.ticker.C:
			// увеличиваем счётчик свободных запросов
			l.count = min(l.maxCount, l.count+1)
		}
	}
}

func (l *limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// есть свободные запросы, получаем управление
	case <-l.ch:
		return nil
	}
}

// создаёт лимитер
// d - частота запросов
// count - количество запросов
func NewLimiter(d time.Duration, count int) *limiter {
	l := &limiter{
		maxCount: count,
		count:    count,
		ticker:   time.NewTicker(d / time.Duration(count)),
		ch:       make(chan struct{}),
	}
	go l.run()

	return l
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
