package kafka

import (
	"encoding/json"
	"time"

  "github.com/JamesAndresCM/eco_orders/internal/orders"
	"github.com/IBM/sarama"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
)

type OrderProcessedEvent struct {
	UserID      int         `json:"user_id"`
	OrderID     int         `json:"order_id"`
	Items       []orders.OrderItem `json:"items"`
	Total       string      `json:"total"`
	ProcessedAt string      `json:"processed_at"`
}

type OrderFailedEvent struct {
	UserID   int         `json:"user_id"`
	Items    []orders.OrderItem `json:"items"`
	Error    string      `json:"error"`
	FailedAt string      `json:"failed_at"`
}

type Producer struct {
	sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V3_6_0_0

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{p}, nil
}

func (p *Producer) PublishOrderFailed(
	userID int,
	items []orders.OrderItem,
	err error,
) {
	if p == nil {
		logger.Warn("kafka producer is nil, skipping publish")
		return
	}

	payload := OrderFailedEvent{
		UserID:   userID,
		Items:    items,
		Error:    err.Error(),
		FailedAt: time.Now().UTC().Format(time.RFC3339),
	}

	bytes, _ := json.Marshal(payload)

	logger.Info("publishing orders.failed event:", string(bytes))

	msg := &sarama.ProducerMessage{
		Topic: "orders.failed",
		Value: sarama.ByteEncoder(bytes),
	}

	if _, _, err := p.SyncProducer.SendMessage(msg); err != nil {
		logger.Error("failed to publish orders.failed:", err)
		return
	}

	logger.Success("orders.failed event published")
}

func (p *Producer) PublishOrderProcessed(userID int, orderID int, items []orders.OrderItem, total string) error {
	payload := OrderProcessedEvent{
		UserID:      userID,
		OrderID:     orderID,
		Items:       items,
		Total:       total,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	}

	bytes, _ := json.Marshal(payload)

	logger.Info("publishing orders.processed event:", string(bytes))

	msg := &sarama.ProducerMessage{
		Topic: "orders.processed",
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err := p.SyncProducer.SendMessage(msg)
	return err
}
