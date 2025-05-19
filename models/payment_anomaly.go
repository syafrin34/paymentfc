package models

import "time"

type PaymentAnomaly struct {
	ID         int    `json:"id"`
	OrderID    int64  `json:"order_id"`
	ExternalID string `json:"external_id"`
	// 1 : Anomaly Amount
	AnomalyType int    `json:"anomaly_type"`
	Notes       string `json:"notes"`
	// 1 : Success
	// 2 : Need to check
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
