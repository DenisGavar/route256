package converter

import (
	"route256/checkout/internal/domain/model"
	desc "route256/checkout/pkg/checkout_v1"
)

func FromDescToMolelAddToCartRequest(addToCartRequest *desc.AddToCartRequest) *model.AddToCartRequest {
	if addToCartRequest == nil {
		return nil
	}

	return &model.AddToCartRequest{
		User:  addToCartRequest.GetUser(),
		Sku:   addToCartRequest.GetSku(),
		Count: addToCartRequest.GetCount(),
	}
}

func FromDescToMolelListCartRequest(listCartRequest *desc.ListCartRequest) *model.ListCartRequest {
	if listCartRequest == nil {
		return nil
	}

	return &model.ListCartRequest{
		User: listCartRequest.GetUser(),
	}
}

func FromMolelToDescListCartResponse(listCartResponse *model.ListCartResponse) *desc.ListCartResponse {
	if listCartResponse == nil {
		return nil
	}

	items := make([]*desc.CartItem, 0, len(listCartResponse.Items))
	for _, i := range listCartResponse.Items {
		items = append(items, FromMolelToDescCartItem(i))
	}

	return &desc.ListCartResponse{
		Items:      items,
		TotalPrice: listCartResponse.TotalPrice,
	}
}

func FromMolelToDescCartItem(cartItem *model.CartItem) *desc.CartItem {
	if cartItem == nil {
		return nil
	}

	return &desc.CartItem{
		Sku:   cartItem.Sku,
		Count: cartItem.Count,
		Name:  cartItem.Name,
		Price: cartItem.Price,
	}
}

func FromDescToMolelPurchaseRequest(purchaseRequest *desc.PurchaseRequest) *model.PurchaseRequest {
	if purchaseRequest == nil {
		return nil
	}

	return &model.PurchaseRequest{
		User: purchaseRequest.User,
	}
}

func FromMolelToDescPurchaseResponse(purchaseResponse *model.PurchaseResponse) *desc.PurchaseResponse {
	if purchaseResponse == nil {
		return nil
	}

	return &desc.PurchaseResponse{
		OrderId: purchaseResponse.OrderId,
	}
}

func FromDescToMolelDeleteFromCartRequest(deleteFromCartRequest *desc.DeleteFromCartRequest) *model.DeleteFromCartRequest {
	if deleteFromCartRequest == nil {
		return nil
	}

	return &model.DeleteFromCartRequest{
		User:  deleteFromCartRequest.User,
		Sku:   deleteFromCartRequest.Sku,
		Count: deleteFromCartRequest.Count,
	}
}
