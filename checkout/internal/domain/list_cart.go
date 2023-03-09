package domain

import (
	"context"
	"route256/checkout/internal/domain/model"
	"route256/checkout/pkg/product-service_v1"

	"github.com/pkg/errors"
)

func (s *service) ListCart(ctx context.Context, req *model.ListCartRequest) (*model.ListCartResponse, error) {
	// получаем список товаров в корзине по user64, в цикле опрашиваем ProductService
	// пока списка товаров нет, идём в ProductService за произвольным товаром

	response, err := s.repository.checkoutRepository.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, cartItem := range response.Items {
		product, err := s.productServiceClient.GetProduct(ctx, &product.GetProductRequest{Sku: cartItem.Sku})
		if err != nil {
			return nil, errors.WithMessage(err, "getting product")
		}

		cartItem.Name = product.Name
		cartItem.Price = product.Price
		response.TotalPrice += product.Price * cartItem.Count
	}

	return response, nil
}
