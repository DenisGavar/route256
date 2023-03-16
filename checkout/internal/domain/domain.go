package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/domain/model"
	product "route256/checkout/pkg/product-service_v1"
	loms "route256/loms/pkg/loms_v1"
)

var (
	ErrNotEnoughItems = errors.New("not enough items")
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type CheckoutRepository interface {
	AddToCart(ctx context.Context, addToCartRequest *model.AddToCartRequest) error
	ListCart(ctx context.Context, listCartRequest *model.ListCartRequest) (*model.ListCartResponse, error)
	DeleteFromCart(ctx context.Context, deleteFromRequest *model.DeleteFromCartRequest) error
}

type repository struct {
	checkoutRepository CheckoutRepository
	transactionManager TransactionManager
}

func NewRepository(checkoutRepository CheckoutRepository, transactionManager TransactionManager) repository {
	return repository{
		checkoutRepository: checkoutRepository,
		transactionManager: transactionManager,
	}
}

var _ Service = (*service)(nil)

type Service interface {
	AddToCart(context.Context, *model.AddToCartRequest) error
	DeleteFromCart(context.Context, *model.DeleteFromCartRequest) error
	ListCart(context.Context, *model.ListCartRequest) (*model.ListCartResponse, error)
	Purchase(context.Context, *model.PurchaseRequest) (*model.PurchaseResponse, error)
}

type ProductServiceClient interface {
	GetProduct(context.Context, *product.GetProductRequest) (*product.GetProductResponse, error)
}

type LomsClient interface {
	Stocks(context.Context, *loms.StocksRequest) (*loms.StocksResponse, error)
	CreateOrder(context.Context, *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error)
}

type service struct {
	lomsClient           LomsClient
	productServiceClient ProductServiceClient
	repository           repository
}

func NewService(lomsClient LomsClient, productServiceClient ProductServiceClient, repository repository) *service {
	return &service{
		lomsClient:           lomsClient,
		productServiceClient: productServiceClient,
		repository:           repository,
	}
}
