package domain

import (
	"context"
	"errors"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	"route256/loms/internal/domain/model"
	repository "route256/loms/internal/repository/postgres"
	repositoryMock "route256/loms/internal/repository/postgres/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrderPayed(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *gomock.Controller) repository.LomsRepository

	type args struct {
		ctx context.Context
		req *model.OrderPayedRequest
	}

	var (
		mc     = gomock.NewController(t)
		ctx    = context.Background()
		dbMock = transactorMock.NewMockDB(mc)

		orderId = gofakeit.Int64()

		req = &model.OrderPayedRequest{
			OrderId: orderId,
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().ClearReserves(ctx, req.OrderId).Return(nil)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusPayed).Return(nil)
				return mock
			},
		},
		{
			name: "clearing reserves fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrClearingReserves,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().ClearReserves(ctx, req.OrderId).Return(repositoryErr)
				return mock
			},
		},
		{
			name: "changing status fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrChangingStatus,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().ClearReserves(ctx, req.OrderId).Return(nil)
				mock.EXPECT().ChangeStatus(ctx, req.OrderId, model.OrderStatusPayed).Return(repositoryErr)
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

			err := client.OrderPayed(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
