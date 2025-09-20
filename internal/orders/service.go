package orders

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JamesAndresCM/eco_orders/pkg/logger"
	"github.com/shopspring/decimal"
)

type Service struct {
	Repo *OrderRepository
}

func (s *Service) ProcessOrder(req OrderRequest) (OrderResponse, error) {
	logger.Info("Processing order:", req)

	if len(req.Items) == 0 {
		return OrderResponse{}, ErrEmptyItems
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	orderID, err := s.Repo.CreateOrder(ctx, req.UserID)
	if err != nil {
		tx.Rollback()
		return OrderResponse{}, fmt.Errorf("failed to create order: %w", err)
	}

	total := decimal.NewFromInt(0)

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			tx.Rollback()
			return OrderResponse{}, ErrInvalidQuantity
		}

		price, stock, err := s.Repo.GetProduct(ctx, item.ProductID)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return OrderResponse{}, ErrProductNotFound
		}
		if err != nil {
			tx.Rollback()
			return OrderResponse{}, fmt.Errorf("failed to fetch product: %w", err)
		}

		if stock < item.Quantity {
			tx.Rollback()
			return OrderResponse{}, ErrInsufficientStock
		}

		if err := s.Repo.InsertOrderItem(ctx, orderID, item.ProductID, item.Quantity, price); err != nil {
			tx.Rollback()
			return OrderResponse{}, fmt.Errorf("failed to insert order item: %w", err)
		}

		if err := s.Repo.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
			tx.Rollback()
			return OrderResponse{}, fmt.Errorf("failed to update product stock: %w", err)
		}

		qty := decimal.NewFromInt(int64(item.Quantity))
		total = total.Add(price.Mul(qty))
	}

	if err := s.Repo.UpdateOrderTotal(ctx, orderID, total); err != nil {
		tx.Rollback()
		return OrderResponse{}, fmt.Errorf("failed to update order total: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return OrderResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	resp := OrderResponse{
		OrderID: orderID,
		Status:  "processed",
		Total:   total,
	}

	return resp, nil
}
