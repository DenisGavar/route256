package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	repositoryMock "route256/checkout/internal/repository/postgres/mocks"
	"route256/libs/transactor"
	transactorMock "route256/libs/transactor/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestDeleteFromCart(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *gomock.Controller) repository.CheckoutRepository
	type dbMockFunc func(mc *gomock.Controller) transactor.DB

	type args struct {
		ctx context.Context
		req *model.DeleteFromCartRequest
	}

	var (
		mc    = gomock.NewController(t)
		ctx   = context.Background()
		tx    = transactorMock.NewMockTx(mc)
		ctxTx = context.WithValue(ctx, transactor.TxKey, tx)
		opts  = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}

		userID    = gofakeit.Int64()
		itemSku   = gofakeit.Uint32()
		itemCount = gofakeit.Uint32()

		req = &model.DeleteFromCartRequest{
			User:  userID,
			Sku:   itemSku,
			Count: itemCount,
		}

		repositoryErr = errors.New("repository error")
	)

	tests := []struct {
		name                   string
		args                   args
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		dbMock                 dbMockFunc
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
				mock.EXPECT().DeleteFromCart(ctxTx, req).Return(nil)
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
			name: "deleting from cart fail",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrDeletingFromCart,
			checkoutRepositoryMock: func(mc *gomock.Controller) repository.CheckoutRepository {
				mock := repositoryMock.NewMockCheckoutRepository(mc)
				mock.EXPECT().DeleteFromCart(ctxTx, req).Return(repositoryErr)
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

			repo := NewRepository(tt.checkoutRepositoryMock(mc), transactor.NewTransactionManager(tt.dbMock(mc)))

			client := NewMockService(repo)

			err := client.DeleteFromCart(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}

		})
	}

}
