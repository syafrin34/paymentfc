package usecase

import (
	"context"
	"paymentfc/cmd/payment/service"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"

	"github.com/sirupsen/logrus"
)

type XenditUsecase interface {
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

type xenditUsecase struct {
	xenditService service.XenditService
}

func NewXenditUseCase(xenditService service.XenditService) XenditUsecase {
	return &xenditUsecase{
		xenditService: xenditService,
	}
}

func (uc *xenditUsecase) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	err := uc.xenditService.CreateInvoice(ctx, param)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("uc.xenditService.CreateInvoice got error : %v", err)
		return err
	}
	return nil
}
