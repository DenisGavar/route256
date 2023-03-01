package domain

import (
	"context"
)

type ListCartRequest struct {
	// user ID
	User int64 `protobuf:"varint,1,opt,name=user,proto3" json:"user,omitempty"`
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

func (m *model) ListCart(ctx context.Context, req *ListCartRequest) (*ListCartResponse, error) {
	// получаем список товаров в корзине по user64, в цикле опрашиваем ProductService
	// пока списка товаров нет, идём в ProductService за произвольным товаром

	return &ListCartResponse{
		Items: []*CartItem{
			{
				Sku:   111,
				Count: 1,
				Name:  "qwe",
				Price: 1,
			},
			{
				Sku:   222,
				Count: 2,
				Name:  "qwe",
				Price: 2,
			},
		},
		TotalPrice: 5,
	}, nil

	// var skus []uint32
	// skus = append(skus, 1076963, 1148162)

	// Items := make([]Item, 0, len(skus))

	// var totalPrice uint32
	// for _, sku := range skus {
	// 	product, err := m.productGetter.GetProduct(ctx, sku)
	// 	if err != nil {
	// 		return nil, 0, errors.WithMessage(err, "getting product")
	// 	}

	// 	item := Item{
	// 		SKU:   sku,
	// 		Count: 1,
	// 		Name:  product.Name,
	// 		Price: product.Price,
	// 	}

	// 	Items = append(Items, item)
	// 	totalPrice += product.Price * uint32(item.Count)
	// }

	// return Items, totalPrice, nil
}
