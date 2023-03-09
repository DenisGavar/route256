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

	var skus []uint32
	skus = append(skus, 1076963, 1148162)

	cartItems := make([]*model.CartItem, 0, len(skus))

	var totalPrice uint32
	for _, sku := range skus {
		product, err := s.productServiceClient.GetProduct(ctx, &product.GetProductRequest{Sku: sku})
		if err != nil {
			return nil, errors.WithMessage(err, "getting product")
		}

		item := &model.CartItem{
			Sku:   sku,
			Count: 1,
			Name:  product.GetName(),
			Price: product.GetPrice(),
		}

		cartItems = append(cartItems, item)
		totalPrice += product.GetPrice() * item.Count
	}

	return &model.ListCartResponse{
		Items:      cartItems,
		TotalPrice: totalPrice,
	}, nil
}
