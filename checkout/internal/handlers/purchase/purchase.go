package purchase

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

type Response struct {
	OrderID int64 `json:"orderID"`
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
	log.Printf("purchase: %+v", req)

	var response Response

	order, err := h.businessLogic.Purchase(ctx, req.User)
	if err != nil {
		return response, err
	}

	response.OrderID = order.OrderID

	return response, nil
}