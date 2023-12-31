package domain

import (
	"context"
	"errors"
	productServiceCachedClient "route256/checkout/internal/clients/cache/product-service"
	"route256/checkout/internal/clients/grpc/loms"
	productServiceGRPCClient "route256/checkout/internal/clients/grpc/product-service"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	"route256/libs/limiter"
	"route256/libs/transactor"
	workerPool "route256/libs/worker-pool"
)

var (
	ErrNotEnoughItems   = errors.New("not enough items")
	ErrGettingListCart  = errors.New("getting list cart")
	ErrGettingProduct   = errors.New("getting product")
	ErrCreatingOrder    = errors.New("creating order")
	ErrDeletingFromCart = errors.New("deleting from cart")
	ErrCheckingStocks   = errors.New("checking stocks")
	ErrAddingToCart     = errors.New("adding to cart")
	ErrLimiter          = errors.New("limiter error")
)

type repo struct {
	checkoutRepository repository.CheckoutRepository
	transactionManager transactor.TransactionManager
}

func NewRepository(checkoutRepository repository.CheckoutRepository, transactionManager transactor.TransactionManager) *repo {
	return &repo{
		checkoutRepository: checkoutRepository,
		transactionManager: transactionManager,
	}
}

var _ Service = (*service)(nil)

type productServiceSettings struct {
	limiter limiter.Limiter
}

func NewProductServiceSettings(limiter limiter.Limiter) *productServiceSettings {
	return &productServiceSettings{
		limiter: limiter,
	}
}

type productService struct {
	productServiceClient       productServiceGRPCClient.ProductServiceClient
	productServiceSettings     *productServiceSettings
	productServiceCachedClient productServiceCachedClient.CachedClient
}

func NewProductService(productServiceClient productServiceGRPCClient.ProductServiceClient, productServiceSettings *productServiceSettings, productServiceCachedClient productServiceCachedClient.CachedClient) *productService {
	return &productService{
		productServiceClient:       productServiceClient,
		productServiceSettings:     productServiceSettings,
		productServiceCachedClient: productServiceCachedClient,
	}
}

type Service interface {
	AddToCart(context.Context, *model.AddToCartRequest) error
	DeleteFromCart(context.Context, *model.DeleteFromCartRequest) error
	ListCart(context.Context, *model.ListCartRequest) (*model.ListCartResponse, error)
	Purchase(context.Context, *model.PurchaseRequest) (*model.PurchaseResponse, error)
}

type service struct {
	lomsClient     loms.LomsClient
	productService *productService
	repository     *repo
	wp             workerPool.Pool[*model.CartItem, *model.CartItem]
}

func NewService(lomsClient loms.LomsClient, productService *productService, repo *repo, workerPool workerPool.Pool[*model.CartItem, *model.CartItem]) *service {
	return &service{
		lomsClient:     lomsClient,
		productService: productService,
		repository:     repo,
		wp:             workerPool,
	}
}

func NewMockService(deps ...interface{}) *service {
	ns := service{}

	for _, v := range deps {
		switch s := v.(type) {
		case loms.LomsClient:
			ns.lomsClient = s
		case *productService:
			ns.productService = s
		case *repo:
			ns.repository = s
		case workerPool.Pool[*model.CartItem, *model.CartItem]:
			ns.wp = s
		}
	}

	return &ns
}
