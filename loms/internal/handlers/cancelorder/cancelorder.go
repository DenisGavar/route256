package cancelorder

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

type Response struct{}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("cancel order: %+v", request)
	return Response{}, nil
}
