package main

import (
	"os"

	"github.com/JamesAndresCM/eco_orders/internal/db"
	"github.com/JamesAndresCM/eco_orders/internal/handlers"
	"github.com/JamesAndresCM/eco_orders/internal/orders"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
  "github.com/JamesAndresCM/eco_orders/internal/kafka"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found, relying on environment variables")
	}

	db.InitDB()
	if db.DB != nil {
		defer db.DB.Close()
	}

	repo := &orders.OrderRepository{DB: db.DB}
	svc := &orders.Service{Repo: repo}
  brokers := kafka.BrokersFromEnv()
	logger.Info("Kafka brokers:", brokers)
	handler := &handlers.OrderHandler{Service: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", handler.CreateOrderHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server running on port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Server failed:", err)
	}
}
