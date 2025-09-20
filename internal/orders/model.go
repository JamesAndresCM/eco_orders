package orders

import "github.com/shopspring/decimal"

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderRequest struct {
	UserID int         `json:"user_id"`
	Items  []OrderItem `json:"items"`
}

type OrderResponse struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
  Total   decimal.Decimal `json:"total"`
}
