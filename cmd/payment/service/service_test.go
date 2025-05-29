package service

import (
	"context"
	mocks "paymentfc/cmd/test_mocks"
	"paymentfc/infrastructure/logger"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CheckPaymentAmountByOrderID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mockPaymentService := mocks.NewMockPaymentService(ctrl)
	mockRepositoryDatabase := mocks.NewMockPaymentDatabase(ctrl)
	mockRepositoryPublisher := mocks.NewMockPaymentEventPublisher(ctrl)

	// expected result
	var expectedAmount float64 = 1000
	mockRepositoryDatabase.EXPECT().
		CheckPaymentAmountByOrderID(context.Background(), int64(1)).
		Return(expectedAmount, nil)

	// actual result
	paymentService := paymentService{
		database:  mockRepositoryDatabase,
		publisher: mockRepositoryPublisher,
	}
	//svc := NewPaymentService(mockRepositoryDatabase, mockRepositoryPublisher)
	actualAmount, actualError := paymentService.CheckPaymentAmountByOrderID(context.Background(), int64(1))
	assert.Equal(t, actualAmount, expectedAmount)
	assert.NoError(t, actualError)

}

func Test_CheckPaymentAmountByOrderID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger.SetupLogger()
	mockPaymentDatabase := mocks.NewMockPaymentDatabase(ctrl)
	mockPaymentEventPublisher := mocks.NewMockPaymentEventPublisher(ctrl)

	// expected result
	expectedAmount := float64(0)

	// mock data
	mockPaymentDatabase.EXPECT().CheckPaymentAmountByOrderID(context.Background(), int64(1)).Return(expectedAmount, assert.AnError)

	paymentService := paymentService{
		database:  mockPaymentDatabase,
		publisher: mockPaymentEventPublisher,
	}

	svc := NewPaymentService(paymentService.database, paymentService.publisher)

	actualAmount, err := svc.CheckPaymentAmountByOrderID(context.Background(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, actualAmount, expectedAmount)

}
