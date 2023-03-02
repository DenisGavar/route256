package domain

import (
	"context"
	"route256/checkout/pkg/product-service_v1"

	"github.com/pkg/errors"
)

type ListCartRequest struct {
	// user ID
	User int64
}

type CartItem struct {
	// stock keeping unit - единица складского учёта
	Sku   uint32
	Count uint32
	// наименование товара
	Name string
	// цена товара
	Price uint32
}

type ListCartResponse struct {
	Items      []*CartItem
	TotalPrice uint32
}

func (m *model) ListCart(ctx context.Context, req *ListCartRequest) (*ListCartResponse, error) {
	// получаем список товаров в корзине по user64, в цикле опрашиваем ProductService
	// пока списка товаров нет, идём в ProductService за произвольным товаром

	var skus []uint32
	skus = append(skus, 1076963, 1148162)

	cartItems := make([]*CartItem, 0, len(skus))

	var totalPrice uint32
	for _, sku := range skus {
		product, err := m.productServiceClient.GetProduct(ctx, &product.GetProductRequest{Sku: sku})
		if err != nil {
			return nil, errors.WithMessage(err, "getting product")
		}

		item := &CartItem{
			Sku:   sku,
			Count: 1,
			Name:  product.GetName(),
			Price: product.GetPrice(),
		}

		cartItems = append(cartItems, item)
		totalPrice += product.GetPrice() * item.Count
	}

	return &ListCartResponse{
		Items:      cartItems,
		TotalPrice: totalPrice,
	}, nil
}
