package domain

import (
	"context"
	"route256/checkout/internal/domain/model"
	product "route256/checkout/pkg/product-service_v1"
	workerpool "route256/libs/worker-pool"
	"sync"

	"github.com/pkg/errors"
)

func (s *service) ListCart(ctx context.Context, req *model.ListCartRequest) (*model.ListCartResponse, error) {
	// получаем список товаров в корзине
	listCart, err := s.repository.checkoutRepository.ListCart(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingListCart.Error())
	}

	// для каждого товара получаем наименование и цену, считаем итого

	// инициализируем worker pool
	// количество worker-ов не больше указанного значения в конфиге
	// функция для worker-а дополняет структуру *model.CartItem
	pool, results := workerpool.New[*model.CartItem, *model.CartItem](
		ctx,
		s.productService.productServiceSettings.listCartWorkersCount,
	)
	pool.Init(ctx)

	// функция для worker-а
	callback := func(cartItem *model.CartItem) *workerpool.Result[*model.CartItem] {
		// вкючаем лимитер, значение передаётся через конфиг
		s.productService.productServiceSettings.limiter.Wait(ctx)
		// получаем данные из product service
		productResp, err := s.productService.productServiceClient.GetProduct(ctx, &product.GetProductRequest{Sku: cartItem.Sku})
		if err != nil {
			// если ошибка при получении данных, то возвращаем ошибку
			return &workerpool.Result[*model.CartItem]{
				Out:   nil,
				Error: errors.WithMessage(err, "getting product"),
			}
		}

		// возвращаем дополненную структуру
		cartItem.Price = productResp.GetPrice()
		cartItem.Name = productResp.GetName()
		return &workerpool.Result[*model.CartItem]{
			Out:   cartItem,
			Error: nil,
		}
	}

	// mutex для записи результата
	var mux sync.Mutex

	// добавляем в worker pool каждую sku для обработки
	// не в горутине, т.к. нам надо гарантированно добавить все sku
	for _, cartItem := range listCart.Items {
		pool.Push(ctx, &workerpool.Job[*model.CartItem, *model.CartItem]{
			Callback: callback,
			Args:     cartItem,
		})
	}

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
			}
			mux.Lock()
			response.Items = append(response.Items, result.Out)
			response.TotalPrice += result.Out.Price * result.Out.Count
			mux.Unlock()
			// отмечаем, что задача выполнена, результат получен
			pool.JobDone()
		}
	}()

	// останавливаем worker pool, передаём аргумент, что ждём завершения всех отправленных задач
	pool.Stop(true)

	// если при обращении к product service были ошибки
	if resultErr != nil {
		return nil, resultErr
	}

	return response, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
