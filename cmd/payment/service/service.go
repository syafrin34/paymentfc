package service

import (
	"context"
	"math"
	"time"

	//"go/constant"
	"paymentfc/cmd/payment/repository"
	"paymentfc/infrastructure/constant"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"

	"github.com/sirupsen/logrus"
)

const (
	maxRetryPublish = 5
)

type PaymentService interface {
	ProcessPaymentSuccess(ctx context.Context, orderID int64) error
	CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)
	SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error
	SavePaymentRequests(ctx context.Context, param models.PaymentRequests) error
	GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error)
}

type paymentService struct {
	database  repository.PaymentDatabase
	publisher repository.PaymentEventPublisher
}

func NewPaymentService(db repository.PaymentDatabase, pb repository.PaymentEventPublisher) PaymentService {
	return &paymentService{
		database:  db,
		publisher: pb,
	}
}

func (s *paymentService) CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error) {
	amount, err := s.database.CheckPaymentAmountByOrderID(ctx, orderID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.checkPaymentAmountByOrderID got error: %v", err)
		return 0, err
	}
	return amount, nil

}
func (s *paymentService) ProcessPaymentSuccess(ctx context.Context, orderID int64) error {

	// validate either order id already paid
	isAlreadyPaid, err := s.database.IsAlreadyPaid(ctx, orderID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.IsAlreadyPaid got error: %v", err)
		return err
	}
	if isAlreadyPaid {
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Infof("[skip - order %d] payment status already paid!", orderID)
		return nil
	}

	// implement retry mechanism

	//publish event kafka
	//retry until  max retry publish
	err = retryPublishPayment(maxRetryPublish, func() error {
		errLogAudit := s.database.InsertAuditLog(ctx, models.PaymentAuditLog{
			OrderID:    orderID,
			Event:      "PublishPaymentSuccess",
			Actor:      "payment",
			CreateTime: time.Now(),
		})
		if errLogAudit != nil {
			logger.Logger.WithFields(logrus.Fields{
				"order_id": orderID,
				"event":    "PublishPaymentSuccess",
				"actor":    "payment",
			}).WithError(errLogAudit).Errorf("s.database.InsertAuditLog() got error: %v", errLogAudit)
		}
		return s.publisher.PublishPaymentSuccess(ctx, orderID)
	})
	//err = s.publisher.PublishPaymentSuccess(ctx, orderID)
	if err != nil {
		// store data to DB --> failed_events
		failedEventsParam := models.FailedEvents{
			OrderID:    orderID,
			FailedType: constant.FailedPublishEventPaymentSuccess,
			Status:     constant.FaildePublishEventNeedToCheck,
			Notes:      err.Error(),
			CreateTime: time.Now(),
		}
		// dead letter table
		errSaveFailedPublish := s.database.SaveFailedPublishEvent(ctx, failedEventsParam)
		if errSaveFailedPublish != nil {
			logger.Logger.WithFields(logrus.Fields{
				"failedEventsParam": failedEventsParam,
			}).WithError(errSaveFailedPublish)
			return errSaveFailedPublish
		}
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.publisher.PublishPaymentSuccess got error: %v", err)
		return err
	}

	//update status db
	err = s.database.MarkPaid(ctx, orderID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.MarkPaid got error")
		return err
	} else {
		errLogAudit := s.database.InsertAuditLog(ctx, models.PaymentAuditLog{
			OrderID:    orderID,
			Event:      "MarkPaid",
			Actor:      "payment",
			CreateTime: time.Now(),
		})
		if errLogAudit != nil {
			logger.Logger.WithFields(logrus.Fields{
				"order_id": orderID,
				"event":    "MarkPaid",
				"actor":    "payment",
			}).WithError(errLogAudit).Errorf("s.database.InsertAuditLog() got error: %v", errLogAudit)
		}

	}
	return nil
}
func (s *paymentService) SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error {
	err := s.database.SavePaymentAnomaly(ctx, param)
	if err != nil {
		return err
	}
	return nil
}
func (s *paymentService) SavePaymentRequests(ctx context.Context, param models.PaymentRequests) error {
	err := s.database.SavePaymentRequests(ctx, param)
	if err != nil {
		return err
	}
	return nil
}
func retryPublishPayment(max int, fn func() error) error {
	var err error
	for i := 0; i <= max; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		//publis event
		//failed -->retry
		//set jeda (2)
		//failed --> retry
		wait := time.Duration(math.Pow(2, float64(i))) * time.Second
		logger.Logger.Printf("Retry: %d, Error: %s. retrying in %d seconds...", i+1, err, wait)
		time.Sleep(wait)

	}
	return err
}
func (s *paymentService) GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error) {
	paymentInfo, err := s.database.GetPaymentInfoByOrderID(ctx, orderID)
	if err != nil {
		return models.Payment{}, err
	}
	return paymentInfo, nil
}
