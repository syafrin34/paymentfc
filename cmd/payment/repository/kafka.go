package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"paymentfc/models"

	"github.com/segmentio/kafka-go"
)

type PaymentEventPublisher interface {
	PublishPaymentSuccess(ctx context.Context, orderID int64) error
	PublishEventPaymentStatus(ctx context.Context, orderID int64, status string, topic string) error
}

type kafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(writer *kafka.Writer) PaymentEventPublisher {
	return &kafkaPublisher{
		writer: writer,
	}
}

// publish payment success
func (k *kafkaPublisher) PublishPaymentSuccess(ctx context.Context, OrderID int64) error {
	payload := map[string]interface{}{
		"order_id": OrderID,
		"status":   "paid",
	}
	data, _ := json.Marshal(payload)
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", OrderID)),
		Value: data,
	})
}

// publish payment failed
func (k *kafkaPublisher) PublishPaymentFailed(ctx context.Context, orderID int64) error {
	payload := map[string]interface{}{
		"order_id": orderID,
		"status":   "failed",
	}

	data, _ := json.Marshal(payload)
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", orderID)),
		Value: data,
	})
}

func (k *kafkaPublisher) PublishEventPaymentStatus(ctx context.Context, orderID int64, status string, topic string) error {
	payload := models.PaymentStatusUpdateEvent{
		OrderID: orderID,
		Status:  status,
	}
	data, _ := json.Marshal(payload)
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", orderID)),
		Topic: topic,
		Value: data,
	})
}
