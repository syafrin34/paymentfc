package usecase

import (
	"context"
	"paymentfc/cmd/payment/service"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type PaymentUseCase interface {
	ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error
}

type paymentUseCase struct {
	service service.PaymentService
}

func NewPaymentUseCase(service service.PaymentService) PaymentUseCase {
	return &paymentUseCase{
		service: service,
	}
}

func (uc *paymentUseCase) ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error {
	switch payload.Status {
	case "PAID":
		// construct external id --> order id
		orderID := extractOrderID(payload.ExternalID)
		//connect ke service layer
		err := uc.service.ProcessPaymentSuccess(ctx, orderID)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"status":      payload.Status,
				"external_id": payload.ExternalID,
			}).Errorf("uc.svc.ProcessPaymentSuccess got error: %v ", err)
			return err
		}
	case "FAILED":

	case "PENDING":

	default:
		logger.Logger.WithFields(logrus.Fields{
			"status": payload.Status,
		}).Infof("[%s] Anomaly Webhook Status: %s", payload.ExternalID, payload.Status)

		//next kita akan buat table baru "payment anomaly"
	}

	return nil
}

func extractOrderID(ExternalID string) int64 {
	idStr := strings.Trim(ExternalID, "order-")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}
