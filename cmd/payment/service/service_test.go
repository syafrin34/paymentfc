package service

import (
	"context"
	mocks "paymentfc/cmd/test_mocks"
	"paymentfc/infrastructure/constant"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"
	"testing"
	"time"

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

// gabungan 2 testing di atas dalam satu test
func Test_CheckPaymentAmountByorderID(t *testing.T) {
	type mockFields struct {
		database  *mocks.MockPaymentDatabase
		publisher *mocks.MockPaymentEventPublisher
	}

	type args struct {
		ctx     context.Context
		orderID int64
	}

	// expected result
	expectedAmountSuccess := float64(10000)
	expectedAmountFailed := float64(0)

	//create list of test case
	tests := []struct {
		name       string
		args       args
		mock       func(mockFields)
		wantResult float64
		wantError  error
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
			mock: func(mf mockFields) {
				mf.database.EXPECT().CheckPaymentAmountByOrderID(context.Background(), int64(1)).Return(expectedAmountSuccess, nil)
			},
			wantResult: expectedAmountSuccess,
			wantError:  nil,
		},
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
			mock: func(mf mockFields) {
				mf.database.EXPECT().CheckPaymentAmountByOrderID(context.Background(), int64(1)).Return(expectedAmountFailed, assert.AnError)
			},
			wantResult: expectedAmountFailed,
			wantError:  assert.AnError,
		},
	}

	// looping test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			logger.SetupLogger()
			mocks := mockFields{
				database:  mocks.NewMockPaymentDatabase(ctrl),
				publisher: mocks.NewMockPaymentEventPublisher(ctrl),
			}
			service := &paymentService{
				database:  mocks.database,
				publisher: mocks.publisher,
			}
			test.mock(mocks)
			// call actual function
			gotAmount, gotError := service.CheckPaymentAmountByOrderID(test.args.ctx, test.args.orderID)
			assert.Equal(t, gotAmount, test.wantResult)
			assert.Equal(t, gotError, test.wantError)
		})
	}

}

func Test_SavePaymentAnomaly(t *testing.T) {
	type mockFields struct {
		database  *mocks.MockPaymentDatabase
		publisher *mocks.MockPaymentEventPublisher
	}

	type args struct {
		ctx   context.Context
		param models.PaymentAnomaly
	}

	mockTime := time.Now()

	tests := []struct {
		name      string
		args      args
		mock      func(mockFields)
		wantError error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				param: models.PaymentAnomaly{
					OrderID:     1,
					ExternalID:  "order-1",
					AnomalyType: constant.AnomalyTypeInvalidAmount,
					Notes:       "",
					Status:      constant.PaymentAnomalyStatusNeedToCheck,
					CreateTime:  mockTime,
				},
			},
			mock: func(mf mockFields) {
				mf.database.EXPECT().SavePaymentAnomaly(context.Background(), models.PaymentAnomaly{
					OrderID:     1,
					ExternalID:  "order-1",
					AnomalyType: constant.AnomalyTypeInvalidAmount,
					Notes:       "",
					Status:      constant.PaymentAnomalyStatusNeedToCheck,
					CreateTime:  mockTime,
				}).Return(nil)
			},
			wantError: nil,
		}, {
			name: "error",
			args: args{
				ctx: context.Background(),
				param: models.PaymentAnomaly{
					OrderID:    999,
					ExternalID: "abc",
				},
			},
			mock: func(mf mockFields) {
				mf.database.EXPECT().SavePaymentAnomaly(context.Background(), models.PaymentAnomaly{
					OrderID:    999,
					ExternalID: "abc",
				}).Return(assert.AnError)
			},
			wantError: assert.AnError,
		},
	}

	// looping for test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mockFields{
				database:  mocks.NewMockPaymentDatabase(ctrl),
				publisher: mocks.NewMockPaymentEventPublisher(ctrl),
			}

			service := &paymentService{
				database:  mock.database,
				publisher: mock.publisher,
			}

			test.mock(mock)
			gotErr := service.SavePaymentAnomaly(test.args.ctx, test.args.param)
			assert.Equal(t, gotErr, test.wantError)
		})
	}
}
