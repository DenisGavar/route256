package converter

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

func ToAddToCartRequestModel(addToCartRequest *desc.AddToCartRequest) *domain.AddToCartRequest {
	if addToCartRequest == nil {
		return nil
	}

	return &domain.AddToCartRequest{
		User:  addToCartRequest.GetUser(),
		Sku:   addToCartRequest.GetSku(),
		Count: addToCartRequest.GetCount(),
	}
}

func ToListCartRequestModel(listCartRequest *desc.ListCartRequest) *domain.ListCartRequest {
	if listCartRequest == nil {
		return nil
	}

	return &domain.ListCartRequest{
		User: listCartRequest.GetUser(),
	}
}

func ToListCartResponseDesc(listCartResponse *domain.ListCartResponse) *desc.ListCartResponse {
	if listCartResponse == nil {
		return nil
	}

	items := make([]*desc.CartItem, 0, len(listCartResponse.Items))
	for _, i := range listCartResponse.Items {
		items = append(items, ToCartItemDesc(i))
	}

	return &desc.ListCartResponse{
		Items:      items,
		TotalPrice: listCartResponse.TotalPrice,
	}
}

func ToCartItemDesc(cartItem *domain.CartItem) *desc.CartItem {
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

func ToPurchaseRequestModel(purchaseRequest *desc.PurchaseRequest) *domain.PurchaseRequest {
	if purchaseRequest == nil {
		return nil
	}

	return &domain.PurchaseRequest{
		User: purchaseRequest.User,
	}
}

func ToPurchaseResponseDesc(purchaseResponse *domain.PurchaseResponse) *desc.PurchaseResponse {
	if purchaseResponse == nil {
		return nil
	}

	return &desc.PurchaseResponse{
		OrderId: purchaseResponse.OrderId,
	}
}
