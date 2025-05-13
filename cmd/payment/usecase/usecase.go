package usecase

import (
	"strconv"
	"strings"
	"toko/ecommerce-msa/PAYMENTFC/cmd/payment/service"
	"toko/ecommerce-msa/PAYMENTFC/infrastructure/logger"
	"toko/ecommerce-msa/PAYMENTFC/models"

	"github.com/sirupsen/logrus"
)

type PaymentUseCase interface {
}

type paymentUseCase struct {
	svc service.PaymentService
}

func NewPaymentUseCase(service service.PaymentService) PaymentUseCase {
	return &paymentUseCase{
		svc: service,
	}
}

func (uc *paymentUseCase) ProcessPaymentWebhook(payload models.XenditWebhookPayload) error {
	switch payload.Status {
	case "PAID":
		// construct external id --> order id
		orderID := extractOrderID(payload.ExternalID)
		//connect ke service layer
		err := uc.svc.ProcessPaymentSuccess(orderID)
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
