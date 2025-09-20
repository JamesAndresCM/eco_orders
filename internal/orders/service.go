package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JamesAndresCM/eco_orders/pkg/logger"
	"github.com/shopspring/decimal"
)

var (
	ErrEmptyItems        = errors.New("items cannot be empty")
	ErrInvalidQuantity   = errors.New("quantity must be greater than 0")
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock for product")
)

type Service struct {
	DB *sql.DB
}

func (s *Service) ProcessOrder(req OrderRequest) (OrderResponse, error) {
	logger.Info("Processing order:", req)
	if len(req.Items) == 0 {
		return OrderResponse{}, ErrEmptyItems
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var orderID int
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO orders (user_id, status, total_amount, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())
		 RETURNING id`,
		req.UserID, 0, 0,
	).Scan(&orderID)

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
		var stock_quantity int
		var price decimal.Decimal

		err := tx.QueryRowContext(
			ctx,
			`SELECT price, stock_quantity FROM products WHERE id = $1`,
			item.ProductID,
		).Scan(&price, &stock_quantity)

		if err == sql.ErrNoRows {
			tx.Rollback()
			return OrderResponse{}, ErrProductNotFound
		}
		if err != nil {
			tx.Rollback()
			return OrderResponse{}, fmt.Errorf("failed to fetch product: %w", err)
		}

		if stock_quantity < item.Quantity {
			tx.Rollback()
			return OrderResponse{}, ErrInsufficientStock
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, unit_price, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
			orderID, item.ProductID, item.Quantity, price,
		)
		if err != nil {
			tx.Rollback()
			return OrderResponse{}, fmt.Errorf("failed to insert order item: %w", err)
		}

		// actualizar stock
		_, err = tx.ExecContext(
			ctx,
			`UPDATE products SET stock_quantity = stock_quantity - $1, updated_at = NOW() WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			tx.Rollback()
			return OrderResponse{}, fmt.Errorf("failed to update product stock: %w", err)
		}

		qty := decimal.NewFromInt(int64(item.Quantity))
		subtotal := price.Mul(qty)
		total = total.Add(subtotal)
	}

	// Actualizar total en la orden
	_, err = tx.ExecContext(
		ctx,
		`UPDATE orders SET total_amount = $1, status = $2, updated_at = NOW() WHERE id = $3`,
		total, 1, orderID,
	)
	if err != nil {
		tx.Rollback()
		return OrderResponse{}, fmt.Errorf("failed to update order total: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return OrderResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Success(fmt.Sprintf("Order %d processed successfully with total %d", orderID, total))

	resp := OrderResponse{
		OrderID: orderID,
		Status:  "processed",
		Total:   total,
	}
	logger.Success("Order processed successfully:", resp)
	return resp, nil
}
