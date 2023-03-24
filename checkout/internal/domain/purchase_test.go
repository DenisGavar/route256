package domain

import (
	"context"
	"errors"
	"testing"

	lomsClient "route256/checkout/internal/clients/grpc/loms"
	lomsClientMock "route256/checkout/internal/clients/grpc/loms/mocks"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	repositoryMock "route256/checkout/internal/repository/postgres/mocks"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	loms "route256/loms/pkg/loms_v1"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestPurchase(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *gomock.Controller) repository.CheckoutRepository
	type dbMockFunc func(mc *gomock.Controller) transactor.DB
	type lomsClientMockFunc func(mc *gomock.Controller) lomsClient.LomsClient

	type args struct {
		ctx context.Context
		req *model.PurchaseRequest
	}

	var (
		mc    = gomock.NewController(t)
		ctx   = context.Background()
		tx    = transactorMock.NewMockTx(mc)
		ctxTx = context.WithValue(ctx, transactor.TxKey, tx)
		opts  = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}

		userID    = gofakeit.Int64()
		orderID   = gofakeit.Int64()
		itemSku   = gofakeit.Uint32()
		itemCount = gofakeit.Uint32()

		req = &model.PurchaseRequest{
			User: userID,
		}

		reqListCart = &model.ListCartRequest{
			User: userID,
		}

		reqCreateOrder = &loms.CreateOrderRequest{
			User: userID,
			Items: []*loms.OrderItem{
				{
					Sku:   itemSku,
					Count: itemCount,
				},
			},
		}

		reqDeleteFromCart = &model.DeleteFromCartRequest{
			User:  userID,
			Sku:   itemSku,
			Count: itemCount,
		}

		res = &model.PurchaseResponse{
			OrderId: orderID,
		}

		resWithoutOrderId = &model.PurchaseResponse{
			OrderId: 0,
		}

		resListCart = &model.ListCartResponse{
			Items: []*model.CartItem{
				{
					Sku:   itemSku,
					Count: itemCount,
				},
			},
		}

		resListCartEmpty = &model.ListCartResponse{
			Items: []*model.CartItem{},
		}

		resCreateOrder = &loms.CreateOrderResponse{
			OrderId: orderID,
		}

		repositoryErr = errors.New("repository error")
		lomsErr       = errors.New("loms error")
	)

	tests := []struct {
		name                   string
		args                   args
		want                   *model.PurchaseResponse
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		dbMock                 dbMockFunc
		lomsClientMock         lomsClientMockFunc
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
				mock.EXPECT().ListCart(ctxTx, reqListCart).Return(resListCart, nil)
				mock.EXPECT().DeleteFromCart(ctxTx, reqDeleteFromCart).Return(nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Commit(ctx).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().CreateOrder(ctxTx, reqCreateOrder).Return(resCreateOrder, nil)
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
				mock.EXPECT().ListCart(ctxTx, reqListCart).Return(nil, repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				return mock
			},
		},
		{
			name: "empty cart",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: resWithoutOrderId,
			err:  nil,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(ctxTx, reqListCart).Return(resListCartEmpty, nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Commit(ctx).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
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
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(ctxTx, reqListCart).Return(resListCart, nil)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().CreateOrder(ctxTx, reqCreateOrder).Return(nil, lomsErr)
				return mock
			},
		},
		{
			name: "deleting from cart fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrDeletingFromCart,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().ListCart(ctxTx, reqListCart).Return(resListCart, nil)
				mock.EXPECT().DeleteFromCart(ctxTx, reqDeleteFromCart).Return(repositoryErr)
				return mock
			},
			dbMock: func(mc *gomock.Controller) transactor.DB {
				mock := transactorMock.NewMockDB(mc)
				mock.EXPECT().BeginTx(ctx, opts).Return(tx, nil)
				tx.EXPECT().Rollback(ctx).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *gomock.Controller) lomsClient.LomsClient {
				mock := lomsClientMock.NewMockLomsClient(mc)
				mock.EXPECT().CreateOrder(ctxTx, reqCreateOrder).Return(resCreateOrder, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepository(tt.checkoutRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)))

			client := NewMockService(repo, tt.lomsClientMock(mc))

			res, err := client.Purchase(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
