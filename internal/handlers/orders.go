package handlers

import (
	"encoding/json"
	"net/http"
  o "github.com/JamesAndresCM/eco_orders/internal/orders"
  "github.com/JamesAndresCM/eco_orders/pkg/logger"
)

type OrderHandler struct {
    Service *o.Service
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
  logger.Info("Received a request to create an order")
	if r.Method != http.MethodPost {
    logger.Warn("Invalid method:", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req o.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    logger.Error("Failed to decode request:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.ProcessOrder(req)
	if err != nil {
    logger.Error("Error processing order:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  logger.Success("Order processed successfully:", resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
