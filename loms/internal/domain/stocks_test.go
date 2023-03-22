package domain

import (
	"context"
	"errors"
	"testing"

	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	"route256/loms/internal/domain/model"
	repository "route256/loms/internal/repository/postgres"
	repositoryMock "route256/loms/internal/repository/postgres/mocks"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestStocks(t *testing.T) {
	type lomsRepositoryMockFunc func(mc *gomock.Controller) repository.LomsRepository

	type args struct {
		ctx context.Context
		req *model.StocksRequest
	}

	var (
		mc     = gomock.NewController(t)
		ctx    = context.Background()
		dbMock = transactorMock.NewMockDB(mc)

		itemSku     = gofakeit.Uint32()
		warehouseId = gofakeit.Int64()
		stocksCount = gofakeit.Uint64()

		req = &model.StocksRequest{
			Sku: itemSku,
		}

		res = &model.StocksResponse{
			Stocks: []*model.StockItem{
				{
					WarehouseId: warehouseId,
					Count:       stocksCount,
				},
			},
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		want               *model.StocksResponse
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().Stocks(ctx, req).Return(res, nil)
				return mock
			},
		},
		{
			name: "getting stocks fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrGettingStocks,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().Stocks(ctx, req).Return(nil, repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepository(tt.lomsRepositoryMock(mc), transactor.NewTransactionManager(dbMock))

			client := NewService(repo)

			res, err := client.Stocks(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
