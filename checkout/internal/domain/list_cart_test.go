package domain

import (
	"context"
	"testing"

	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	repositoryMock "route256/checkout/internal/repository/postgres/mocks"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	type checkoutRepositoryMockFunc func(mc *gomock.Controller) repository.CheckoutRepository

	type args struct {
		ctx context.Context
		req *model.ListCartRequest
	}

	var (
		mc  = gomock.NewController(t)
		ctx = context.Background()

		userID = gofakeit.Int64()

		req = &model.ListCartRequest{
			User: userID,
		}

		resRepository = &model.ListCartResponse{
			Items: []*model.CartItem{
				&model.CartItem{
					Sku:   123,
					Count: 2,
				},
			},
		}

		res = &model.ListCartResponse{
			Items: []*model.CartItem{
				&model.CartItem{
					Sku:   123,
					Count: 2,
					Name:  "123",
					Price: 123,
				},
			},
			TotalPrice: 246,
		}
	)

	tests := []struct {
		name                   string
		args                   args
		want                   *model.ListCartResponse
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
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
				mock.EXPECT().ListCart(ctx, req).Return(resRepository)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := NewMockService(tt.checkoutRepositoryMock(mc))

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
