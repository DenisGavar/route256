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
	// канал с результатами выполнения задач
	results chan *Result[Out]

	// WaitGroup для задач в работе
	inWorkWG *sync.WaitGroup
	// WaitGroup для задач в очереди
	pushedWG *sync.WaitGroup
	// признак остановки работы
	stopped bool
}

// конструктор
func New[In, Out any](ctx context.Context, workersCount int) (*pool[In, Out], <-chan *Result[Out]) {
	pool := &pool[In, Out]{
		workersCount: workersCount,
		jobs:         make(chan *Job[In, Out], workersCount),
		results:      make(chan *Result[Out], workersCount),
		inWorkWG:     new(sync.WaitGroup),
		pushedWG:     new(sync.WaitGroup),
		stopped:      false,
	}

	return pool, pool.results
}

// запускаем pool worker-ов
func (p *pool[In, Out]) Init(ctx context.Context) {
	// итерируемся по количеству worker-ов
	for i := 0; i < p.workersCount; i++ {
		// инициализируем каждого worker-а
		go func() {
			p.initWorker(ctx, p.jobs, p.results)
		}()
	}
}

// инициализируем каждого worker
func (p *pool[In, Out]) initWorker(ctx context.Context, jobs <-chan *Job[In, Out], result chan<- *Result[Out]) {
	// читаем задачи из канала
	for job := range jobs {
		select {
		case <-ctx.Done():
			return
		case result <- job.Callback(job.Args):
			// отмечаем выполнение задачи
			p.inWorkWG.Done()
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
			// отмечаем, что задача добавлена на исполнение
			p.inWorkWG.Add(1)
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
	// ждём завершения всех задач в работе
	p.inWorkWG.Wait()
	// закрываем канал с результатами работы задач
	close(p.results)
}

// отмечаем, что задача выполнена, результат обработан
func (p *pool[In, Out]) JobDone() {
	p.pushedWG.Done()
}
