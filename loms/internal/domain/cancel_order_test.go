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

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *gomock.Controller) repository.LomsRepository
	type dbMockFunc func(mc *gomock.Controller) transactor.DB

	type args struct {
		ctx context.Context
		req *model.CancelOrderRequest
	}

	var (
		mc  = gomock.NewController(t)
		ctx = context.Background()
		tx  = transactorMock.NewMockTx(mc)
		// TODO: use this context instead of gomock.Any()
		//ctxTx = context.WithValue(ctx, transactor.TxKey, tx)
		opts = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}

		orderId     = gofakeit.Int64()
		itemSku     = gofakeit.Uint32()
		itemCount   = gofakeit.Uint64()
		warehouseId = gofakeit.Int64()

		req = &model.CancelOrderRequest{
			OrderId: orderId,
		}

		reqReturnReserve = &model.ReserveStocksItem{
			WarehouseId: warehouseId,
			Sku:         itemSku,
			Count:       itemCount,
		}

		resReserves = &model.Reserve{
			ReserveItems: []*model.ReserveStocksItem{
				{
					WarehouseId: warehouseId,
					Sku:         itemSku,
					Count:       itemCount,
				},
			},
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
		dbMock             dbMockFunc
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
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().Reserves(gomock.Any(), orderId).Return(resReserves, nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ReturnReserve(gomock.Any(), reqReturnReserve).Return(nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ClearReserves(gomock.Any(), orderId).Return(nil)
				mock.EXPECT().ChangeStatus(gomock.Any(), orderId, model.OrderStatusCancelled).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(gomock.Any(), opts).Return(tx, nil)
				tx.EXPECT().Commit(gomock.Any()).Return(nil)
				return mock
			},
		},
		{
			name: "checking reserves fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrCheckingReserves,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().Reserves(gomock.Any(), orderId).Return(nil, repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(gomock.Any(), opts).Return(tx, nil)
				tx.EXPECT().Rollback(gomock.Any()).Return(nil)
				return mock
			},
		},
		{
			name: "returning reserves fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrReturningReserves,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().Reserves(gomock.Any(), orderId).Return(resReserves, nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ReturnReserve(gomock.Any(), reqReturnReserve).Return(repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(gomock.Any(), opts).Return(tx, nil)
				tx.EXPECT().Rollback(gomock.Any()).Return(nil)
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
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().Reserves(gomock.Any(), orderId).Return(resReserves, nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ReturnReserve(gomock.Any(), reqReturnReserve).Return(nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ClearReserves(gomock.Any(), orderId).Return(repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(gomock.Any(), opts).Return(tx, nil)
				tx.EXPECT().Rollback(gomock.Any()).Return(nil)
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
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().Reserves(gomock.Any(), orderId).Return(resReserves, nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ReturnReserve(gomock.Any(), reqReturnReserve).Return(nil)
				// TODO: use ctxTx instead of gomock.Any()
				mock.EXPECT().ClearReserves(gomock.Any(), orderId).Return(nil)
				mock.EXPECT().ChangeStatus(gomock.Any(), orderId, model.OrderStatusCancelled).Return(repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(gomock.Any(), opts).Return(tx, nil)
				tx.EXPECT().Commit(gomock.Any()).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger.Init()

			repo := NewRepository(tt.lomsRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)))

			client := NewService(repo)

			err := client.CancelOrder(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
