package list_order

import (
	"context"
	"errors"
	"log"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

var (
	ErrEmptyOrderID = errors.New("empty order ID")
)

func (r Request) Validate() error {
	if r.OrderID == 0 {
		return ErrEmptyOrderID
	}
	return nil
}

type Response struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Item `json:"items"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("list order: %+v", request)
	return Response{
		Status: "new",
		User:   42,
		Items: []Item{
			{
				SKU:   1,
				Count: 15,
			},
			{
				SKU:   2,
				Count: 33,
			},
		},
	}, nil
}
