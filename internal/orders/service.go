package orders

import (
  "errors"
  "github.com/JamesAndresCM/eco_orders/pkg/logger"
)

var nextOrderID = 1

func ProcessOrder(req OrderRequest) (OrderResponse, error) {
  logger.Info("Processing order:", req)
	if len(req.Items) == 0 {
		return OrderResponse{}, errors.New("order must contain at least one item")
	}

	total := 0
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return OrderResponse{}, errors.New("quantity must be greater than 0")
		}
		total += item.Quantity * 100 
	}

	orderID := nextOrderID
	nextOrderID++

	resp := OrderResponse{
		OrderID: orderID,
		Status:  "processed",
		Total:   total,
	}
  logger.Success("Order processed successfully:", resp)
	return resp, nil
}
