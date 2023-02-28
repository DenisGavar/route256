package domain

import (
	"context"

	"github.com/pkg/errors"
)

type Product struct {
	Name  string
	Price uint32
}

type Item struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

func (m *Model) ListCart(ctx context.Context, user int64) ([]Item, uint32, error) {
	// получаем список товаров в корзине по user64, в цикле опрашиваем ProductService
	// пока списка товаров нет, идём в ProductService за произвольным товаром

	var skus []uint32
	skus = append(skus, 1076963, 1148162)

	Items := make([]Item, 0, len(skus))

	var totalPrice uint32
	for _, sku := range skus {
		product, err := m.productGetter.GetProduct(ctx, sku)
		if err != nil {
			return nil, 0, errors.WithMessage(err, "getting product")
		}

		item := Item{
			SKU:   sku,
			Count: 1,
			Name:  product.Name,
			Price: product.Price,
		}

		Items = append(Items, item)
		totalPrice += product.Price * uint32(item.Count)
	}

	return Items, totalPrice, nil
}
