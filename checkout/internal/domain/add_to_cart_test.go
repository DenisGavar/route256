package domain

import (
	"context"
	"errors"
	lomsClient "route256/checkout/internal/clients/grpc/loms"
	lomsClientMock "route256/checkout/internal/clients/grpc/loms/mocks"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	repositoryMock "route256/checkout/internal/repository/postgres/mocks"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	loms "route256/loms/pkg/loms_v1"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddToCart(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *gomock.Controller) repository.CheckoutRepository
	type lomsClientMockFunc func(mc *gomock.Controller) lomsClient.LomsClient

	type args struct {
		ctx context.Context
		req *model.AddToCartRequest
	}

	var (
		mc     = gomock.NewController(t)
		ctx    = context.Background()
		dbMock = transactorMock.NewMockDB(mc)

		userID               = gofakeit.Int64()
		itemSku              = gofakeit.Uint32()
		itemCount            = gofakeit.Uint32()
		warehouseId          = gofakeit.Int64()
		stocksCount          = uint64(itemCount) + 1
		stocksCountNotEnough = uint64(itemCount) - 1

		req = &model.AddToCartRequest{
			User:  userID,
			Sku:   itemSku,
			Count: itemCount,
		}

		reqStocks = &loms.StocksRequest{
			Sku: itemSku,
		}

		resStocks = &loms.StocksResponse{
			Stocks: []*loms.StockItem{
				{
					WarehouseId: warehouseId,
					Count:       stocksCount,
				},
			},
		}

		resStocksNotEnough = &loms.StocksResponse{
			Stocks: []*loms.StockItem{
				{
					WarehouseId: warehouseId,
					Count:       stocksCountNotEnough,
				},
			},
		}

		repositoryErr = errors.New("repository error")
		lomsErr       = errors.New("loms error")
	)

	tests := []struct {
		name                   string
		args                   args
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		lomsClientMock         lomsClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().AddToCart(ctx, req).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().Stocks(ctx, reqStocks).Return(resStocks, nil)
				return mock
			},
		},
		{
			name: "checking stocks fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrCheckingStocks,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().Stocks(ctx, reqStocks).Return(nil, lomsErr)
				return mock
			},
		},
		{
			name: "adding to cart fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrAddingToCart,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().AddToCart(ctx, req).Return(repositoryErr)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().Stocks(ctx, reqStocks).Return(resStocks, nil)
				return mock
			},
		},
		{
			name: "not enough items fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrNotEnoughItems,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().Stocks(ctx, reqStocks).Return(resStocksNotEnough, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepository(tt.checkoutRepositoryMock(mc), transactor.NewTransactionManager(dbMock))

			client := NewMockService(repo, tt.lomsClientMock(mc))

			err := client.AddToCart(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}

}
