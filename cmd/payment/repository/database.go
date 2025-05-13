package repository

import (
	"context"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentDatabase interface {
	MarkPaid(ctx context.Context, orderID int64) error
}

type paymentDatabase struct {
	DB *gorm.DB
}

func NewPaymentDatabase(db *gorm.DB) PaymentDatabase {
	return &paymentDatabase{
		DB: db,
	}
}

func (r *paymentDatabase) MarkPaid(ctx context.Context, ordeID int64) error {
	// update status db menjadi paid

	err := r.DB.Model(&models.Payment{}).Table("payments").WithContext(ctx).Where("order_id = ?", ordeID).Update("status", "paid").Error
	if err != nil {
		{
			logger.Logger.WithFields(logrus.Fields{
				"order_id": ordeID,
			}).Errorf("r.DB.Update got error: %v ", err)
			return err
		}
	}
	return nil
}
