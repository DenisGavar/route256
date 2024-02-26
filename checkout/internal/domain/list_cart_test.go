package domain

import (
	"context"
	"errors"
	cachedProductService "route256/checkout/internal/clients/cache/product-service"
	productServiceGRPCClient "route256/checkout/internal/clients/grpc/product-service"
	productServiceGRPCClientMock "route256/checkout/internal/clients/grpc/product-service/mocks"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	repositoryMock "route256/checkout/internal/repository/postgres/mocks"
	product "route256/checkout/pkg/product-service_v1"
	"route256/libs/cache"
	"route256/libs/limiter"
	limiterMock "route256/libs/limiter/mocks"
	"route256/libs/logger"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	workerPool "route256/libs/worker-pool"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *gomock.Controller) repository.CheckoutRepository
	type limiterMockFunc func(mc *gomock.Controller) limiter.Limiter
	type productServiceClientMockFunc func(mc *gomock.Controller) productServiceGRPCClient.ProductServiceClient

	type args struct {
		ctx context.Context
		req *model.ListCartRequest
	}

	var (
		mc     = gomock.NewController(t)
		ctx    = context.Background()
		dbMock = transactorMock.NewMockDB(mc)

		listCartWorkersCount = gofakeit.IntRange(5, 10)

		userID    = gofakeit.Int64()
		itemSku   = gofakeit.Uint32()
		itemCount = gofakeit.Uint32()
		itemName  = gofakeit.BeerName()
		itemPrice = gofakeit.Uint32()

		req = &model.ListCartRequest{
			User: userID,
		}

		reqGetProduct = &product.GetProductRequest{
			Sku: itemSku,
		}

		resListCart = &model.ListCartResponse{
			Items: []*model.CartItem{
				{
					Sku:   itemSku,
					Count: itemCount,
				},
			},
		}

		resGetProduct = &product.GetProductResponse{
			Name:  itemName,
			Price: itemPrice,
		}

		res = &model.ListCartResponse{
			Items: []*model.CartItem{
				{
					Sku:   itemSku,
					Count: itemCount,
					Name:  itemName,
					Price: itemPrice,
				},
			},
			TotalPrice: itemPrice * itemCount,
		}

		repositoryErr     = errors.New("repository error")
		productServiceErr = errors.New("product service error")
	)

	tests := []struct {
		name                     string
		args                     args
		want                     *model.ListCartResponse
		err                      error
		checkoutRepositoryMock   checkoutRepositoryMockFunc
		limiterMock              limiterMockFunc
		productServiceClientMock productServiceClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(gomock.Any(), req).Return(resListCart, nil)
				return mock
			},
			limiterMock: func(mc *gomock.Controller) limiter.Limiter {
				mock := limiterMock.NewMockLimiter(mc)
				mock.EXPECT().Wait(gomock.Any()).Return(nil)
				return mock
			},
			productServiceClientMock: func(mc *gomock.Controller) productServiceGRPCClient.ProductServiceClient {
				mock := productServiceGRPCClientMock.NewMockProductServiceClient(mc)
				mock.EXPECT().GetProduct(gomock.Any(), reqGetProduct).Return(resGetProduct, nil)
				return mock
			},
		},
		{
			name: "getting list cart fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrGettingListCart,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(gomock.Any(), req).Return(nil, repositoryErr)
				return mock
			},
			limiterMock: func(mc *gomock.Controller) limiter.Limiter {
				mock := limiterMock.NewMockLimiter(mc)
				return mock
			},
			productServiceClientMock: func(mc *gomock.Controller) productServiceGRPCClient.ProductServiceClient {
				mock := productServiceGRPCClientMock.NewMockProductServiceClient(mc)
				return mock
			},
		},
		{
			name: "getting product fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrGettingProduct,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(gomock.Any(), req).Return(resListCart, nil)
				return mock
			},
			limiterMock: func(mc *gomock.Controller) limiter.Limiter {
				mock := limiterMock.NewMockLimiter(mc)
				mock.EXPECT().Wait(gomock.Any()).Return(nil)
				return mock
			},
			productServiceClientMock: func(mc *gomock.Controller) productServiceGRPCClient.ProductServiceClient {
				mock := productServiceGRPCClientMock.NewMockProductServiceClient(mc)
				mock.EXPECT().GetProduct(gomock.Any(), reqGetProduct).Return(nil, productServiceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger.Init()

			wp := workerPool.New[*model.CartItem, *model.CartItem](
				tt.args.ctx,
				listCartWorkersCount,
			)
			wp.Init(tt.args.ctx)

			repo := NewRepository(tt.checkoutRepositoryMock(mc), transactor.NewTransactionManager(dbMock))

			lruCache := cache.NewCache(10, time.Second*10)

			productServiceCachedClient := cachedProductService.NewCachedClient(lruCache)

			productServiceSettings := NewProductServiceSettings(tt.limiterMock(mc))
			productService := NewProductService(tt.productServiceClientMock(mc), productServiceSettings, productServiceCachedClient)

			client := NewMockService(repo, productService, wp)

			res, err := client.ListCart(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

			wp.Stop(false)

		})
	}

}
