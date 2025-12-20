package kafka

import (
	"context"
	"log"
  "github.com/JamesAndresCM/eco_orders/pkg/logger"
	"github.com/IBM/sarama"
)

func StartOrderConsumer(brokers []string) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(
		brokers,
		"orders-service",
		config,
	)
	if err != nil {
		logger.Error("❌ failed to create consumer group:", err)
	}

	consumer := &OrderConsumer{}
	ctx := context.Background()

	for {
		if err := consumerGroup.Consume(
			ctx,
			[]string{"orders.create"},
			consumer,
		); err != nil {
			logger.Error("❌ consume error:", err)
		}
	}
}

