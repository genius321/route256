package service

import (
	context "context"
	"errors"
	checkout "route256/checkout/internal/pkg/checkout"
	loms "route256/checkout/internal/pkg/loms"
	schema "route256/checkout/internal/repository/schema"
	"testing"

	"github.com/brianvoe/gofakeit"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func Test_Purchase(t *testing.T) {
	t.Parallel()

	t.Run("error while validate req", func(t *testing.T) {
		t.Parallel()

		// Act
		_, err := (&Service{}).Purchase(context.Background(), &checkout.PurchaseRequest{User: -6})

		// Assert
		require.Error(t, err)
	})

	t.Run("error while list cart from repository", func(t *testing.T) {
		t.Parallel()

		errStub := errors.New("stub")

		repositoryMock := NewMockRepository(t)
		repositoryMock.On("ListCart", mock.Anything, mock.Anything).Return(nil, errStub).Once()

		// Act
		_, err := (&Service{
			Repository: repositoryMock,
		}).Purchase(context.Background(), &checkout.PurchaseRequest{User: 5})

		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("error CreateOrder", func(t *testing.T) {
		t.Parallel()

		repositoryMock := NewMockRepository(t)
		transactionManagerMock := NewMockTransactionManager(t)
		lomsClientMock := NewMockLomsClient(t)

		listCartResponse := []*schema.Item{
			{
				Sku:    gofakeit.Int64(),
				Amount: gofakeit.Int64(),
			},
			{
				Sku:    gofakeit.Int64(),
				Amount: gofakeit.Int64(),
			},
		}
		repositoryMock.On("ListCart", mock.Anything, mock.Anything).Return(listCartResponse, nil).Once()

		createOrderResponse := &loms.CreateOrderResponse{OrderId: gofakeit.Int64()}
		errStub := errors.New("stub")
		lomsClientMock.On("CreateOrder", mock.Anything, mock.Anything).Return(createOrderResponse, errStub).Once()

		// deleteAllFromCartResponse := &emptypb.Empty{}
		// repositoryMock.On("DeleteAllFromCart", mock.Anything, mock.Anything).Return(deleteAllFromCartResponse, nil).Once()

		// Act
		_, err := (&Service{
			Repository:         repositoryMock,
			TransactionManager: transactionManagerMock,
			LomsClient:         lomsClientMock,
		}).Purchase(context.Background(), &checkout.PurchaseRequest{User: 5})

		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("error DeleteAllFromCart", func(t *testing.T) {
		t.Parallel()

		repositoryMock := NewMockRepository(t)
		transactionManagerMock := NewMockTransactionManager(t)
		lomsClientMock := NewMockLomsClient(t)

		listCartResponse := []*schema.Item{}
		repositoryMock.On("ListCart", mock.Anything, mock.Anything).Return(listCartResponse, nil).Once()

		createOrderResponse := &loms.CreateOrderResponse{OrderId: gofakeit.Int64()}
		lomsClientMock.On("CreateOrder", mock.Anything, mock.Anything).Return(createOrderResponse, nil).Once()

		errStub := errors.New("stub")
		deleteAllFromCartResponse := &emptypb.Empty{}
		repositoryMock.On("DeleteAllFromCart", mock.Anything, mock.Anything).Return(deleteAllFromCartResponse, errStub).Once()

		// Act
		_, err := (&Service{
			Repository:         repositoryMock,
			TransactionManager: transactionManagerMock,
			LomsClient:         lomsClientMock,
		}).Purchase(context.Background(), &checkout.PurchaseRequest{User: 5})

		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("success Purchase", func(t *testing.T) {
		t.Parallel()

		repositoryMock := NewMockRepository(t)
		transactionManagerMock := NewMockTransactionManager(t)
		lomsClientMock := NewMockLomsClient(t)

		listCartResponse := []*schema.Item{}
		repositoryMock.On("ListCart", mock.Anything, mock.Anything).Return(listCartResponse, nil).Once()

		createOrderResponse := &loms.CreateOrderResponse{OrderId: gofakeit.Int64()}
		lomsClientMock.On("CreateOrder", mock.Anything, mock.Anything).Return(createOrderResponse, nil).Once()

		deleteAllFromCartResponse := &emptypb.Empty{}
		repositoryMock.On("DeleteAllFromCart", mock.Anything, mock.Anything).Return(deleteAllFromCartResponse, nil).Once()

		// Act
		_, err := (&Service{
			Repository:         repositoryMock,
			TransactionManager: transactionManagerMock,
			LomsClient:         lomsClientMock,
		}).Purchase(context.Background(), &checkout.PurchaseRequest{User: 5})

		// Assert
		require.NoError(t, err)
	})
}
