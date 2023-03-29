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

func TestListOrder(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *gomock.Controller) repository.LomsRepository

	type args struct {
		ctx context.Context
		req *model.ListOrderRequest
	}

	var (
		mc     = gomock.NewController(t)
		ctx    = context.Background()
		dbMock = transactorMock.NewMockDB(mc)

		orderId   = gofakeit.Int64()
		user      = gofakeit.Int64()
		itemSku   = gofakeit.Uint32()
		itemCount = gofakeit.Uint16()

		req = &model.ListOrderRequest{
			OrderId: orderId,
		}

		res = &model.ListOrderResponse{
			Status: model.OrderStatusPayed,
			User:   user,
			Items: []*model.OrderItem{
				{
					Sku:   itemSku,
					Count: itemCount,
				},
			},
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		want               *model.ListOrderResponse
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
				mock.EXPECT().ListOrder(ctx, req).Return(res, nil)
				return mock
			},
		},
		{
			name: "getting list order fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrGettingListOrder,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().ListOrder(ctx, req).Return(nil, repositoryErr)
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

			res, err := client.ListOrder(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
