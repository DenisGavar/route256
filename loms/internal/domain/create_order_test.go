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
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *gomock.Controller) repository.LomsRepository
	type dbMockFunc func(mc *gomock.Controller) transactor.DB

	type args struct {
		ctx context.Context
		req *model.CreateOrderRequest
	}

	var (
		mc    = gomock.NewController(t)
		ctx   = context.Background()
		tx    = transactorMock.NewMockTx(mc)
		ctxTx = context.WithValue(ctx, transactor.TxKey, tx)
		opts  = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}

		user                 = gofakeit.Int64()
		orderId              = gofakeit.Int64()
		itemSku              = gofakeit.Uint32()
		itemCount            = gofakeit.Uint16()
		warehouseId          = gofakeit.Int64()
		warehouseIdMany      = gofakeit.Int64()
		stocksCount          = uint64(itemCount) + 1
		stocksCountEqual     = uint64(itemCount)
		stocksCountNotEnough = uint64(itemCount) - 1

		req = &model.CreateOrderRequest{
			User:    user,
			OrderId: orderId,
			Items: []*model.OrderItem{
				{
					Sku:   itemSku,
					Count: itemCount,
				},
			},
		}

		reqStocks = &model.StocksRequest{
			Sku: itemSku,
		}

		reqReserveItems = &model.ReserveStocksItem{
			WarehouseId: warehouseId,
			Sku:         itemSku,
			Count:       uint64(itemCount),
		}

		res = &model.CreateOrderResponse{
			OrderId: orderId,
		}

		resStocks = &model.StocksResponse{
			Stocks: []*model.StockItem{
				{
					WarehouseId: warehouseId,
					Count:       stocksCount,
				},
			},
		}

		resStocksManyWarehouses = &model.StocksResponse{
			Stocks: []*model.StockItem{
				{
					WarehouseId: warehouseId,
					Count:       stocksCountEqual,
				},
				{
					WarehouseId: warehouseIdMany,
					Count:       stocksCount,
				},
			},
		}

		resStocksNotEnough = &model.StocksResponse{
			Stocks: []*model.StockItem{
				{
					WarehouseId: warehouseId,
					Count:       stocksCountNotEnough,
				},
			},
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		want               *model.CreateOrderResponse
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
			want: res,
			err:  nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(resStocks, nil)
				mock.EXPECT().ReserveItems(ctxTx, orderId, reqReserveItems).Return(nil)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusAwaitingPayment).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Commit(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "creating order fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrCreatingOrder,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(nil, repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				return mock
			},
		},
		{
			name: "getting stocks fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(nil, repositoryErr)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusFailed).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "not enough items fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(resStocksNotEnough, nil)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusFailed).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "positive case (stocks count equal)",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(resStocksManyWarehouses, nil)
				mock.EXPECT().ReserveItems(ctxTx, orderId, reqReserveItems).Return(nil)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusAwaitingPayment).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Commit(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "reserving items fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(resStocks, nil)
				mock.EXPECT().ReserveItems(ctxTx, orderId, reqReserveItems).Return(repositoryErr)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusFailed).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "changing status fail (awaiting payment)",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrChangingStatus,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(resStocks, nil)
				mock.EXPECT().ReserveItems(ctxTx, orderId, reqReserveItems).Return(nil)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusAwaitingPayment).Return(repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Commit(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "changing status fail (failed)",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrChangingStatus,
			lomsRepositoryMock: func(mc *gomock.Controller) repository.LomsRepository {
				mock := repositoryMock.NewMockLomsRepository(mc)
				mock.EXPECT().CreateOrder(ctx, req).Return(res, nil)
				mock.EXPECT().Stocks(ctxTx, reqStocks).Return(resStocksNotEnough, nil)
				mock.EXPECT().ChangeStatus(ctx, orderId, model.OrderStatusFailed).Return(repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepository(tt.lomsRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)))

			client := NewService(repo)

			res, err := client.CreateOrder(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
