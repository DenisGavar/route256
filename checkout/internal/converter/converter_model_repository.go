package converter

import (
	"route256/checkout/internal/domain/model"
	"route256/checkout/internal/repository/schema"
)

func FromRepositoryToMolelListCartResponse(cartItems []*schema.CartItem) *model.ListCartResponse {
	if cartItems == nil {
		return nil
	}

	items := make([]*model.CartItem, 0, len(cartItems))
	for _, i := range cartItems {
		items = append(items, FromRepositoryToMolelCartItem(i))
	}

	return &model.ListCartResponse{
		Items: items,
	}
}

func FromRepositoryToMolelCartItem(cartItem *schema.CartItem) *model.CartItem {
	return &model.CartItem{
		Sku:   cartItem.Sku,
		Count: cartItem.Count,
	}
}
