package domain

import (
	"context"
	"errors"
	"testing"

	productServiceGRPCClient "route256/checkout/internal/clients/grpc/product-service"
	productServiceGRPCClientMock "route256/checkout/internal/clients/grpc/product-service/mocks"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	repositoryMock "route256/checkout/internal/repository/postgres/mocks"
	product "route256/checkout/pkg/product-service_v1"
	"route256/libs/limiter"
	limiterMock "route256/libs/limiter/mocks"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	type checkoutRepositoryMockFunc func(mc *gomock.Controller) repository.CheckoutRepository
	type transactionManagerMockFunc func(mc *gomock.Controller) transactor.TransactionManager
	type limiterMockFunc func(mc *gomock.Controller) limiter.Limiter
	type productServiceClientMockFunc func(mc *gomock.Controller) productServiceGRPCClient.ProductServiceClient

	type args struct {
		ctx context.Context
		req *model.ListCartRequest
	}

	var (
		mc  = gomock.NewController(t)
		ctx = context.Background()

		listCartWorkersCount = gofakeit.IntRange(5, 10)

		userID    = gofakeit.Int64()
		itemSku   = gofakeit.Uint32()
		itemCount = gofakeit.Uint32()
		itemName  = gofakeit.BeerName()
		itemPrice = gofakeit.Uint32()

		req = &model.ListCartRequest{
			User: userID,
		}

		reqProductService = &product.GetProductRequest{
			Sku: itemSku,
		}

		resRepository = &model.ListCartResponse{
			Items: []*model.CartItem{
				{
					Sku:   itemSku,
					Count: itemCount,
				},
			},
		}

		resProductService = &product.GetProductResponse{
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

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name                     string
		args                     args
		want                     *model.ListCartResponse
		err                      error
		checkoutRepositoryMock   checkoutRepositoryMockFunc
		transactionManagerMock   transactionManagerMockFunc
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
				mock.EXPECT().ListCart(ctx, req).Return(resRepository, nil)
				return mock
			},
			transactionManagerMock: func(mc *gomock.Controller) transactor.TransactionManager {
				mock := transactorMock.NewMockTransactionManager(mc)
				return mock
			},
			limiterMock: func(mc *gomock.Controller) limiter.Limiter {
				mock := limiterMock.NewMockLimiter(mc)
				mock.EXPECT().Wait(ctx).Return(nil)
				return mock
			},
			productServiceClientMock: func(mc *gomock.Controller) productServiceGRPCClient.ProductServiceClient {
				mock := productServiceGRPCClientMock.NewMockProductServiceClient(mc)
				mock.EXPECT().GetProduct(ctx, reqProductService).Return(resProductService, nil)
				return mock
			},
		},
		{
			name: "checkoutRepository.ListCart fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrGettingListCart,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(ctx, req).Return(nil, repositoryErr)
				return mock
			},
			transactionManagerMock: func(mc *gomock.Controller) transactor.TransactionManager {
				mock := transactorMock.NewMockTransactionManager(mc)
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepository(tt.checkoutRepositoryMock(mc), tt.transactionManagerMock(mc))

			productServiceSettings := NewProductServiceSettings(listCartWorkersCount, tt.limiterMock(mc))
			productService := NewProductService(tt.productServiceClientMock(mc), productServiceSettings)

			client := NewMockService(repo, productService)

			res, err := client.ListCart(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
