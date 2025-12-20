package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/JamesAndresCM/eco_orders/internal/orders"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
)

func StartOrderConsumer(brokers []string, orderService *orders.Service, producer *Producer) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(
		brokers,
		"orders-service",
		config,
	)
	if err != nil {
		logger.Error("failed to create consumer group:", err)
		return
	}

	consumer := &OrderConsumer{
		OrderService: orderService,
		Producer:     producer,
	}
	ctx := context.Background()

	for {
		if err := consumerGroup.Consume(
			ctx,
			[]string{"orders.create"},
			consumer,
		); err != nil {
			logger.Error("consume error:", err)
		}
	}
}
