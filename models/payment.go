package models

import (
	"time"
)

type Payment struct {
	ID          int64     `json:"id"`
	OrderID     int64     `json:"order_id"`
	UserID      int64     `json:"user_id"`
	ExternalID  string    `json:"external_id"`
	Amount      float64   `json:"amlount"`
	Status      string    `json:"status"`
	CreateTime  time.Time `json:"create_time"`
	ExpiredTime time.Time `json:"expired_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type PaymentRequests struct {
	ID         int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
	Notes      string    `json:"notes"`
	CreateTime time.Time `json:"create_time"`
	UpdateYime time.Time `json:"update_time"`
}

type PaymentStatusUpdateEvent struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}
type FailedPaymentList struct {
	TotalFailedPayment int `json:"total_failed_payments"`
	PaymentList        []PaymentRequests
}
