package main

import (
	"log"
	"net/http"
  o "github.com/JamesAndresCM/eco_orders/internal/handlers"
  "github.com/JamesAndresCM/eco_orders/pkg/logger"
);

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", o.CreateOrderHandler)

  logger.Logger.Println("Starting eco-orders service on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
