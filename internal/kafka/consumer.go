package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/JamesAndresCM/eco_orders/internal/orders"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
)

type OrderCreateEvent struct {
	UserID int `json:"user_id"`
	Items  []struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

type OrderConsumer struct {
	OrderService *orders.Service
	Producer     *Producer
}

func (c *OrderConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *OrderConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *OrderConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {
		var event OrderCreateEvent

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("invalid message:", err)
			session.MarkMessage(msg, "")
			continue
		}

		logger.Info("processing order from kafka user_id=", event.UserID)

		req := orders.OrderRequest{
			UserID: event.UserID,
			Items:  make([]orders.OrderItem, 0, len(event.Items)),
		}

		for _, item := range event.Items {
			req.Items = append(req.Items, orders.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		resp, err := c.OrderService.ProcessOrder(req)
		if err != nil {
			logger.Error("failed to process order:", err)

			if c.Producer != nil {
				c.Producer.PublishOrderFailed(
					event.UserID,
					req.Items,
					err,
				)
			} else {
				logger.Error("kafka producer is nil, cannot publish orders.failed")
			}

			session.MarkMessage(msg, "")
			continue
		}

		if c.Producer != nil {
			if err := c.Producer.PublishOrderProcessed(
				event.UserID,
				resp.OrderID,
				req.Items,
				resp.Total.String(),
			); err != nil {
				logger.Error("failed to publish orders.processed:", err)
			}
		}

		logger.Success("order processed successfully user_id=", event.UserID)
		session.MarkMessage(msg, "")
	}

	return nil
}

