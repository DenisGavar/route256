package listcart

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
)

type Handler struct {
	businessLogic *domain.Model
}

func New(businessLogic *domain.Model) *Handler {
	return &Handler{
		businessLogic: businessLogic,
	}
}

type Request struct {
	User int64 `json:"user"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Response struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"totalPrice"`
}

var (
	ErrEmptyUser = errors.New("empty user")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	var response Response

	items, totalPrice, err := h.businessLogic.ListCart(ctx, req.User)
	if err != nil {
		return response, err
	}

	response.Items = make([]Item, 0, len(items))
	for _, item := range items {
		response.Items = append(response.Items, Item{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}
	response.TotalPrice = totalPrice

	return response, nil
}
