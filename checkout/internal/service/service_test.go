package service

import (
	"context"
	"errors"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"
	"route256/checkout/internal/repository/schema"
	"route256/checkout/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_Purchase(t *testing.T) {
	errStub := errors.New("stub")
	repositoryMock := mocks.NewRepositoryMock(t)
	lomsClientMock := mocks.NewLomsClientMock(t)

	type listCartResp struct {
		resp []*schema.Item
		err  error
	}

	type lomsCreateOrderResp struct {
		resp *loms.CreateOrderResponse
		err  error
	}

	type deleteAllFromCartResp struct {
		resp *emptypb.Empty
		err  error
	}

	type args struct {
		ctx                   context.Context
		req                   *checkout.PurchaseRequest
		listCartResp          *listCartResp
		lomsCreateOrderResp   *lomsCreateOrderResp
		deleteAllFromCartResp *deleteAllFromCartResp
	}

	tests := []struct {
		name    string
		args    args
		want    *checkout.PurchaseResponse
		wantErr bool
	}{
		{
			name: "should be error, req.User == -6",
			args: args{
				ctx: context.Background(),
				req: &checkout.PurchaseRequest{User: -6},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should be error, repository.listCart error",
			args: args{
				ctx: context.Background(),
				req: &checkout.PurchaseRequest{User: 5},
				listCartResp: &listCartResp{
					resp: nil,
					err:  errStub,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should be error, lomsClientMock.CreateOrderMock error",
			args: args{
				ctx: context.Background(),
				req: &checkout.PurchaseRequest{User: 5},
				listCartResp: &listCartResp{
					resp: nil,
					err:  nil,
				},
				lomsCreateOrderResp: &lomsCreateOrderResp{
					resp: nil,
					err:  errStub,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should be error, repositoryMock.DeleteAllFromCartMock error",
			args: args{
				ctx: context.Background(),
				req: &checkout.PurchaseRequest{User: 5},
				listCartResp: &listCartResp{
					resp: nil,
					err:  nil,
				},
				lomsCreateOrderResp: &lomsCreateOrderResp{
					resp: &loms.CreateOrderResponse{OrderId: 1},
					err:  nil,
				},
				deleteAllFromCartResp: &deleteAllFromCartResp{
					resp: nil,
					err:  errStub,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should be ok, Purchase",
			args: args{
				ctx: context.Background(),
				req: &checkout.PurchaseRequest{User: 5},
				listCartResp: &listCartResp{
					resp: []*schema.Item{
						{
							Sku:    1,
							Amount: 2,
						},
						{
							Sku:    3,
							Amount: 4,
						},
					},
					err: nil,
				},
				lomsCreateOrderResp: &lomsCreateOrderResp{
					resp: &loms.CreateOrderResponse{OrderId: 1},
					err:  nil,
				},
				deleteAllFromCartResp: &deleteAllFromCartResp{
					resp: nil,
					err:  nil,
				},
			},
			want: &checkout.PurchaseResponse{OrderId: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			if tt.args.listCartResp != nil {
				repositoryMock.ListCartMock.Return(
					tt.args.listCartResp.resp,
					tt.args.listCartResp.err,
				)
			}
			if tt.args.lomsCreateOrderResp != nil {
				lomsClientMock.CreateOrderMock.Return(
					tt.args.lomsCreateOrderResp.resp,
					tt.args.lomsCreateOrderResp.err,
				)
			}
			if tt.args.deleteAllFromCartResp != nil {
				repositoryMock.DeleteAllFromCartMock.Return(
					tt.args.deleteAllFromCartResp.resp,
					tt.args.deleteAllFromCartResp.err,
				)
			}
			result, err := (&Service{
				Repository: repositoryMock,
				LomsClient: lomsClientMock,
			}).Purchase(tt.args.ctx, tt.args.req)
			// Assert
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, result)
		})
	}
}
