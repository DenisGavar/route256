package workerpool

import (
	"context"
	"sync"
)

// сруктура задач на выполнение
type Job[In, Out any] struct {
	// некая функция, которую нам необходмо выполнить
	Callback func(In) *Result[Out]
	// аргументы функции, которую необходимо выполнить
	Args In
	// канал с результатами выполнения задач
	Results chan *Result[Out]
}

// структура результата
type Result[Out any] struct {
	// возвращаем результат выполнения функции
	Out Out
	// возвращаем ошибку, если есть
	Error error
}

// интерфейс worker pool
type Pool[In, Out any] interface {
	// инициализация worker pool, запуск worker-ов
	Init(ctx context.Context)
	// добавление задачи
	Push(ctx context.Context, job *Job[In, Out])
	// остановка worker pool
	Stop(wait bool)
	// пометка выполенной и обработанной задачи
	JobDone()
}

// для проверки удовлетворения интерфейсу
var _ Pool[any, any] = &pool[any, any]{}

// worker pool
type pool[In, Out any] struct {
	// количество worker
	workersCount int

	// канал с задачами на выполнение
	jobs chan *Job[In, Out]

	// WaitGroup для задач в очереди
	pushedWG *sync.WaitGroup
	// признак остановки работы
	stopped bool
}

// конструктор
func New[In, Out any](ctx context.Context, workersCount int) *pool[In, Out] {
	pool := &pool[In, Out]{
		workersCount: workersCount,
		jobs:         make(chan *Job[In, Out], workersCount),
		pushedWG:     new(sync.WaitGroup),
		stopped:      false,
	}

	return pool
}

// запускаем pool worker-ов
func (p *pool[In, Out]) Init(ctx context.Context) {
	// итерируемся по количеству worker-ов
	for i := 0; i < p.workersCount; i++ {
		// инициализируем каждого worker-а
		go func() {
			p.initWorker(ctx, p.jobs)
		}()
	}
}

// инициализируем каждого worker
func (p *pool[In, Out]) initWorker(ctx context.Context, jobs <-chan *Job[In, Out]) {
	// читаем задачи из канала
	for job := range jobs {
		select {
		case <-ctx.Done():
			return
		case job.Results <- job.Callback(job.Args):
		}

	}
}

// добавляем задачу на обработку
func (p *pool[In, Out]) Push(ctx context.Context, job *Job[In, Out]) {
	// если pool worker остановлен, то новую задачу на обработку не добавляем
	if p.stopped {
		return
	}

	// отмечаем, что задача встала в очередь
	p.pushedWG.Add(1)

	// ждём когда можно будет добавить задачу на выполнение
	go func() {
		select {
		case <-ctx.Done():
			return
		case p.jobs <- job:
		}
	}()
}

// останавливаем pool worker
func (p *pool[In, Out]) Stop(wait bool) {
	// останавливаем pool worker
	p.stopped = true
	if wait {
		p.pushedWG.Wait()
	}
	// закрываем канал для входящих задач
	close(p.jobs)

}

// отмечаем, что задача выполнена, результат обработан
func (p *pool[In, Out]) JobDone() {
	p.pushedWG.Done()
}
