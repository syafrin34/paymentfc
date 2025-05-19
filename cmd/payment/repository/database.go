package repository

import (
	"context"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentDatabase interface {
	CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)
	MarkPaid(ctx context.Context, orderID int64) error
	SavePayment(ctx context.Context, param models.Payment) error
	IsAlreadyPaid(ctx context.Context, orderID int64) (bool, error)
	SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error
}

type paymentDatabase struct {
	DB *gorm.DB
}

func NewPaymentDatabase(db *gorm.DB) PaymentDatabase {
	return &paymentDatabase{
		DB: db,
	}
}
func (r *paymentDatabase) CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error) {
	var result models.Payment
	err := r.DB.Table("payments").WithContext(ctx).Where("order_id=?", orderID).First(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Amount, nil
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
func (r *paymentDatabase) IsAlreadyPaid(ctx context.Context, orderID int64) (bool, error) {
	var result models.Payment
	err := r.DB.Table("payments").WithContext(ctx).Where("external_id=?", orderID).First(&result).Error
	if err != nil {
		return false, err
	}
	return result.Status == "PAID", nil
}
func (r *paymentDatabase) SavePayment(ctx context.Context, param models.Payment) error {
	err := r.DB.Table("payments").WithContext(ctx).Create(param).Error

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("r.DB.Create got error: %v ", err)
		return err
	}
	return nil
}
func (r *paymentDatabase) SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error {
	err := r.DB.Table("payment_anomalies").WithContext(ctx).Create(param).Error
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("r.DB.Create.payment_anomalies got error: %v ", err)
		return err
	}
	return nil

}
