package kafka

import (
	"os"
	"strings"
)

func BrokersFromEnv() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"} // fallback local
	}

	return strings.Split(brokers, ",")
}
