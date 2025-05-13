package service

import (
	"toko/ecommerce-msa/PAYMENTFC/cmd/payment/repository"
	"toko/ecommerce-msa/PAYMENTFC/infrastructure/logger"

	"github.com/sirupsen/logrus"
)

type PaymentService interface {
	ProcessPaymentSuccess(orderID int64) error
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
func (s *paymentService) ProcessPaymentSuccess(orderID int64) error {
	//publish event kafka
	err := s.publisher.PublishPaymentSuccess(orderID)
	if err != nil {

		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.publisher.PublishPaymentSuccess got error: %v", err)
		return err
	}

	//update status db
	err = s.database.MarkPaid(orderID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.MarkPaid got error")
		return err
	}
	return nil
}
