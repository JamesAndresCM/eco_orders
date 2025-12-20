package kafka

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
)

type OrderFailedEvent struct {
	UserID   int         `json:"user_id"`
	Items    interface{} `json:"items"`
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
	items interface{},
	err error,
) {
	if p == nil {
		logger.Warn("⚠️ kafka producer is nil, skipping publish")
		return
	}

	payload := OrderFailedEvent{
		UserID:   userID,
		Items:    items,
		Error:    err.Error(),
		FailedAt: time.Now().UTC().Format(time.RFC3339),
	}

	bytes, _ := json.Marshal(payload)

	logger.Info("📤 publishing orders.failed event:", string(bytes))

	msg := &sarama.ProducerMessage{
		Topic: "orders.failed",
		Value: sarama.ByteEncoder(bytes),
	}

	if _, _, err := p.SyncProducer.SendMessage(msg); err != nil {
		logger.Error("❌ failed to publish orders.failed:", err)
		return
	}

	logger.Success("✅ orders.failed event published")
}
