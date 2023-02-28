package deletefromcart

import (
	"context"
	"errors"
	"log"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct{}

var (
	ErrEmptyUser  = errors.New("empty user")
	ErrEmptySKU   = errors.New("empty sku")
	ErrEmptyCount = errors.New("empty count")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	if r.SKU == 0 {
		return ErrEmptySKU
	}
	if r.Count == 0 {
		return ErrEmptyCount
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("deleteFromCart: %+v", req)

	return Response{}, nil
}
