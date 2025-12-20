package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
)

type OrderCreateEvent struct {
	UserID int `json:"user_id"`
	Items  []struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

type OrderConsumer struct{}

func (c *OrderConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *OrderConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *OrderConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {
		var event OrderCreateEvent

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("❌ invalid message:", err)
			session.MarkMessage(msg, "")
			continue
		}

		logger.Success(
			"✅ order received",
			"user_id=", event.UserID,
			"items=", len(event.Items),
		)

		// 👉 acá después llamamos a ProcessOrder(event)

		session.MarkMessage(msg, "")
	}

	return nil
}
