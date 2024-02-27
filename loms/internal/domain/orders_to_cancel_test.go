package domain

import (
	"context"
	"errors"
	"route256/libs/logger"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	"route256/loms/internal/domain/model"
	repository "route256/loms/internal/repository/postgres"
	repositoryMock "route256/loms/internal/repository/postgres/mocks"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrdersToCancel(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *gomock.Controller) repository.LomsRepository

	type args struct {
		ctx context.Context
		req time.Time
	}

	var (
		mc     = gomock.NewController(t)
		ctx    = context.Background()
		dbMock = transactorMock.NewMockDB(mc)
		req    = time.Now()

		orderID = gofakeit.Int64()

		res = []*model.CancelOrderRequest{
			{
				OrderId: orderID,
			},
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		want               []*model.CancelOrderRequest
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
				mock.EXPECT().OrdersToCancel(gomock.Any(), req).Return(res, nil)
				return mock
			},
		},
		{
			name: "getting orders to cancel fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrGettingOrdersToCancel,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().OrdersToCancel(gomock.Any(), req).Return(nil, repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger.Init()

			repo := NewRepository(tt.lomsRepositoryMock(mc), transactor.NewTransactionManager(dbMock))

			client := NewService(repo)

			res, err := client.OrdersToCancel(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
