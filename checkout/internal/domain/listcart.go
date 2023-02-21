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

func (m *Model) ListCart(ctx context.Context, user int64) ([]Item, error) {
	// получаем список товаров в корзине по user64, в цикле опрашиваем ProductService
	// пока списка товаров нет, идём в ProductService за произвольным товаром

	var sku uint32
	sku = 1076963
	product, err := m.productGetter.GetProduct(ctx, sku)
	if err != nil {
		return nil, errors.WithMessage(err, "getting product")
	}

	return []Item{
		{
			SKU:   sku,
			Count: 42,
			Name:  product.Name,
			Price: product.Price,
		},
	}, nil
}
