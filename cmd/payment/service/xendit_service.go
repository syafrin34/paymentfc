package service

import (
	"context"
	"fmt"
	"paymentfc/cmd/payment/repository"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"
	"time"

	"github.com/sirupsen/logrus"
)

type XenditService interface {
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

type xenditService struct {
	database repository.PaymentDatabase
	xendit   repository.XenditClient
}

func NewXenditService(database repository.PaymentDatabase, xenditClinet repository.XenditClient) XenditService {
	return &xenditService{
		database: database,
		xendit:   xenditClinet,
	}
}
func (s *xenditService) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	//construct external id --> ini kita yang define sendiri

	externalID := fmt.Sprintf("order-%d", param.OrderID)
	req := models.XenditInvoiceRequest{
		ExternalID:  externalID,
		Amount:      param.TotalAmount,
		Description: fmt.Sprintf("[FC] Pembayaran Order %d", param.OrderID),
		PayerEmail:  fmt.Sprintf("user%d@test.com", param.UserID),
	}

	xenditInvoiceDetail, err := s.xendit.CreateInvoice(ctx, req)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"rquest":  param,
			"payload": req,
		}).Errorf("s.xenditCreateInvoice got error : %v", err)
		return err
	}

	//save payment to db

	newPayment := models.Payment{
		OrderID:     param.OrderID,
		UserID:      param.UserID,
		ExternalID:  externalID,
		Amount:      param.TotalAmount,
		Status:      "PENDING", // sweeping status "PENDING"
		CreateTime:  time.Now(),
		ExpiredTime: xenditInvoiceDetail.ExpiryDate,
	}

	err = s.database.SavePayment(ctx, newPayment)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"param":      param,
			"newPayment": newPayment,
		}).Errorf("s.database.SavePayment got error : %v", err)
		return err
	}
	return nil

}
