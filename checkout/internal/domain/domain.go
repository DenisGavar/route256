package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/domain/model"
	product "route256/checkout/pkg/product-service_v1"
	workerPool "route256/libs/worker-pool"
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

type Limiter interface {
	Wait(context.Context) error
}

type productServiceSettings struct {
	limiter Limiter
}

func NewProductServiceSettings(limiter Limiter) *productServiceSettings {
	return &productServiceSettings{
		limiter: limiter,
	}
}

type productService struct {
	productServiceClient   ProductServiceClient
	productServiceSettings productServiceSettings
}

func NewProductService(productServiceClient ProductServiceClient, productServiceSettings productServiceSettings) productService {
	return productService{
		productServiceClient:   productServiceClient,
		productServiceSettings: productServiceSettings,
	}
}

type LomsClient interface {
	Stocks(context.Context, *loms.StocksRequest) (*loms.StocksResponse, error)
	CreateOrder(context.Context, *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error)
}

type service struct {
	lomsClient     LomsClient
	productService productService
	repository     repository
	wp             workerPool.Pool[*model.CartItem, *model.CartItem]
}

func NewService(lomsClient LomsClient, productService productService, repository repository, workerPool workerPool.Pool[*model.CartItem, *model.CartItem]) *service {
	return &service{
		lomsClient:     lomsClient,
		productService: productService,
		repository:     repository,
		wp:             workerPool,
	}
}
