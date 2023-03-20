package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/clients/grpc/loms"
	productServiceGRPCClient "route256/checkout/internal/clients/grpc/product-service"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	"route256/libs/limiter"
	"route256/libs/transactor"
)

var (
	ErrNotEnoughItems = errors.New("not enough items")
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
	listCartWorkersCount int
	limiter              limiter.Limiter
}

func NewProductServiceSettings(listCartWorkersCount int, limiter limiter.Limiter) *productServiceSettings {
	return &productServiceSettings{
		listCartWorkersCount: listCartWorkersCount,
		limiter:              limiter,
	}
}

type productService struct {
	productServiceClient   productServiceGRPCClient.ProductServiceClient
	productServiceSettings productServiceSettings
}

func NewProductService(productServiceClient productServiceGRPCClient.ProductServiceClient, productServiceSettings productServiceSettings) productService {
	return productService{
		productServiceClient:   productServiceClient,
		productServiceSettings: productServiceSettings,
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
	productService productService
	repository     *repo
}

func NewService(lomsClient loms.LomsClient, productService productService, repo *repo) *service {
	return &service{
		lomsClient:     lomsClient,
		productService: productService,
		repository:     repo,
	}
}
