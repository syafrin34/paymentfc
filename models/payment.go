package models

import "time"

type Payment struct {
	ID         int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	ExternalID string    `json:"external_id"`
	Amount     float64   `json:"amlount"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
}
