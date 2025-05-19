package usecase

import (
	"context"
	"errors"
	"fmt"
	"paymentfc/cmd/payment/service"
	"paymentfc/infrastructure/constant"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"
	"strconv"
	"strings"
	"time"

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
		//validate webhook amount before payment success
		amount, err := uc.service.CheckPaymentAmountByOrderID(ctx, orderID)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"order_id":    orderID,
				"status":      payload.Status,
				"external_id": payload.ExternalID,
			})
		}
		if amount != payload.Amount {
			// insert to table payment anomalies
			error_message := fmt.Sprintf("webhook amount mismatch: expected %.2f, got %.2f ", amount, payload.Amount)
			paymentAnomaly := models.PaymentAnomaly{
				OrderID:     orderID,
				ExternalID:  payload.ExternalID,
				AnomalyType: constant.AnomalyTypeInvalidAmount,
				Notes:       error_message,
				Status:      constant.PaymentAnomalyStatusNeedToCheck,
				CreateTime:  time.Now(),
			}
			err := uc.service.SavePaymentAnomaly(ctx, paymentAnomaly)
			if err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"payload":        payload,
					"paymentAnomaly": paymentAnomaly,
				}).WithError(err)
				return err
			}

			logger.Logger.WithFields(logrus.Fields{
				"payload": payload,
			}).Errorf("Webhook amount mismatch: expected %.2f, got %.2f", amount, payload.Amount)
			err = errors.New(error_message)
			return err
			// abort process
			// insert to payment anomaly for future manual checking
		}
		//connect ke service layer
		err = uc.service.ProcessPaymentSuccess(ctx, orderID)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"status":         payload.Status,
				"external_id":    payload.ExternalID,
				"webhook_amount": payload.Amount,
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
