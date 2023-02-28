package create_order

import (
	"context"
	"errors"
	"log"
)

type Request struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

var (
	ErrEmptyUser  = errors.New("empty user")
	ErrEmptyItems = errors.New("empty items")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	if len(r.Items) == 0 {
		return ErrEmptyItems
	}
	return nil
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("create order: %+v", request)
	return Response{
		OrderID: 42,
	}, nil
}
