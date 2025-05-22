package repository

import (
	"context"
	"paymentfc/infrastructure/logger"
	"paymentfc/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentDatabase interface {
	CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)
	MarkPaid(ctx context.Context, orderID int64) error
	SavePayment(ctx context.Context, param models.Payment) error
	IsAlreadyPaid(ctx context.Context, orderID int64) (bool, error)
	SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error
	GetPendingInvoices(ctx context.Context) ([]models.Payment, error)
	SavePaymentRequests(ctx context.Context, param models.PaymentRequests) error
	UpdateSuccessPaymentRequests(ctx context.Context, paymentRequestID int64) error
	UpdateFailedPaymentSuccess(ctx context.Context, paymentRequestID int64, notes string) error
	GetPendingPaymentRequests(ctx context.Context, paymentRequests *[]models.PaymentRequests) error
	GetFailedPaymentRequests(ctx context.Context, paymentRequests *[]models.PaymentRequests) error
	UpdatePendingPaymentRequests(ctx context.Context, paymentRequestID int64) error
	GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error)
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
func (r *paymentDatabase) GetPendingInvoices(ctx context.Context) ([]models.Payment, error) {
	var result []models.Payment
	// data di db > 10mil data
	err := r.DB.Table("payments").WithContext(ctx).Where("status=? AND create_time >= now() - interval '1 day'", "PENDING").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *paymentDatabase) SavePaymentRequests(ctx context.Context, param models.PaymentRequests) error {
	err := r.DB.Table("payment_requests").WithContext(ctx).Create(models.PaymentRequests{
		OrderID:    param.ID,
		UserID:     param.UserID,
		Amount:     param.Amount,
		Status:     param.Status,
		CreateTime: param.CreateTime,
	}).Error
	if err != nil {
		return err
	}

	return nil

}
func (r *paymentDatabase) GetPendingPaymentRequests(ctx context.Context, paymentRequests *[]models.PaymentRequests) error {
	err := r.DB.Table("payment_requests").WithContext(ctx).Where("status=?", "PENDING").Limit(5).Order("create_time ASC").Find(paymentRequests).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentDatabase) GetFailedPaymentRequests(ctx context.Context, paymentRequests *[]models.PaymentRequests) error {
	err := r.DB.Table("payment_requests").WithContext(ctx).Where("status = ?", "FAILED").
		Where("retry_count <= ?", 3).Limit(5).Order("create_time ASC").Find(paymentRequests).Error //

	if err != nil {
		return err
	}
	return nil
}

func (r *paymentDatabase) UpdateSuccessPaymentRequests(ctx context.Context, paymentRequestID int64) error {
	err := r.DB.Table("payment_requests").WithContext(ctx).Where("id=?", paymentRequestID).Updates(map[string]interface{}{
		"status":      "SUCCESS",
		"update_time": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentDatabase) UpdateFailedPaymentSuccess(ctx context.Context, paymentRequestID int64, notes string) error {
	err := r.DB.Table("payment_requests").WithContext(ctx).Where("id=?", paymentRequestID).Updates(map[string]interface{}{
		"status":      "FAILED",
		"notes":       notes,
		"retry_count": gorm.Expr("retry_count + 1"),
		"update_time": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentDatabase) UpdatePendingPaymentRequests(ctx context.Context, paymentRequestID int64) error {
	err := r.DB.Table("payment_requests").WithContext(ctx).Where("id = ? ", paymentRequestID).
		Updates(map[string]interface{}{
			"status":      "PENDING",
			"update_time": time.Now(),
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentDatabase) GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error) {
	var result models.Payment
	err := r.DB.Table("payments").WithContext(ctx).Where("orde_id = ?", orderID).First(&result).Error
	if err != nil {
		return models.Payment{}, err
	}
	return result, nil
}
