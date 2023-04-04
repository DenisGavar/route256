package domain

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	product "route256/checkout/pkg/product-service_v1"
	"route256/libs/logger"
	workerPool "route256/libs/worker-pool"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) ListCart(ctx context.Context, req *model.ListCartRequest) (*model.ListCartResponse, error) {
	// получаем список товаров в корзине
	logger.Debug("checkout domain", zap.String("handler", "ListCart"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout domain ListCart processing")
	defer span.Finish()

	span.SetTag("user", req.User)

	listCart, err := s.repository.checkoutRepository.ListCart(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingListCart.Error())
	}

	// для каждого товара получаем наименование и цену, считаем итого

	// функция для worker-а
	callback := func(cartItem *model.CartItem) *workerPool.Result[*model.CartItem] {
		// вкючаем лимитер, значение передаётся через конфиг
		err := s.productService.productServiceSettings.limiter.Wait(ctx)
		if err != nil {
			return &workerPool.Result[*model.CartItem]{
				Out:   nil,
				Error: errors.WithMessage(err, ErrLimiter.Error()),
			}
		}
		// получаем данные из product service
		productResp, err := s.productService.productServiceClient.GetProduct(ctx, &product.GetProductRequest{Sku: cartItem.Sku})
		if err != nil {
			// если ошибка при получении данных, то возвращаем ошибку
			return &workerPool.Result[*model.CartItem]{
				Out:   nil,
				Error: errors.WithMessage(err, ErrGettingProduct.Error()),
			}
		}

		// возвращаем дополненную структуру
		cartItem.Price = productResp.GetPrice()
		cartItem.Name = productResp.GetName()
		return &workerPool.Result[*model.CartItem]{
			Out:   cartItem,
			Error: nil,
		}
	}

	// mutex для записи результата
	var mux sync.Mutex
	// wg для гарантированного выполнения всех задач
	var wg sync.WaitGroup

	// создаём канал для чтения результатов, буфер на количество запрашиваемых sku
	results := make(chan *workerPool.Result[*model.CartItem], len(listCart.Items))
	// закрываем канал результатов
	defer close(results)

	// добавляем в wg количество задач для обработки
	wg.Add(len(listCart.Items))
	// добавляем в worker pool каждую sku для обработки
	go func() {
		for _, cartItem := range listCart.Items {
			s.wp.Push(ctx, &workerPool.Job[*model.CartItem, *model.CartItem]{
				Callback: callback,
				Args:     cartItem,
				Results:  results,
			})
		}
	}()

	// собираем результат
	var resultErr error
	var response *model.ListCartResponse = &model.ListCartResponse{
		Items: make([]*model.CartItem, 0, len(listCart.Items)),
	}

	go func() {
		// читаем из канала пока он открыт
		for result := range results {
			if result.Error != nil {
				resultErr = result.Error
			} else {
				mux.Lock()
				response.Items = append(response.Items, result.Out)
				response.TotalPrice += result.Out.Price * result.Out.Count
				mux.Unlock()
			}
			// отмечаем, что задача выполнена, результат получен
			s.wp.JobDone()
			// результат получен
			wg.Done()
		}
	}()

	// дожидаемся выполнения всех задач
	wg.Wait()

	// если при обращении к product service были ошибки
	if resultErr != nil {
		return nil, resultErr
	}

	return response, nil
}
