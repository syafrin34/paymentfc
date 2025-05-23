package models

import "time"

type PaymentAuditLog struct {
	ID         int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	PaymentID  int64     `json:"payment_id"`
	ExternalID string    `json:"external_id"`
	Event      string    `json:"event"` // save payment, create invoice, save payment anomaly
	Actor      string    `json:"actor"` // order, xendit, scheduler
	CreateTime time.Time `json:"create_time"`
}
