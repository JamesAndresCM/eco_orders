package orders

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type OrderRepository struct {
	DB *sql.DB
}

func (r *OrderRepository) CreateOrder(ctx context.Context, userID int) (int, error) {
	var orderID int
	err := r.DB.QueryRowContext(
		ctx,
		`INSERT INTO orders (user_id, status, total_amount, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())
		 RETURNING id`,
		userID, 0, 0,
	).Scan(&orderID)
	return orderID, err
}

func (r *OrderRepository) GetProduct(ctx context.Context, productID int) (decimal.Decimal, int, error) {
	var price decimal.Decimal
	var stock int
	err := r.DB.QueryRowContext(
		ctx,
		`SELECT price, stock_quantity FROM products WHERE id = $1`,
		productID,
	).Scan(&price, &stock)
	return price, stock, err
}

func (r *OrderRepository) InsertOrderItem(ctx context.Context, orderID, productID, quantity int, price decimal.Decimal) error {
	_, err := r.DB.ExecContext(
		ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, unit_price, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		orderID, productID, quantity, price,
	)
	return err
}

func (r *OrderRepository) UpdateStock(ctx context.Context, productID, quantity int) error {
	_, err := r.DB.ExecContext(
		ctx,
		`UPDATE products SET stock_quantity = stock_quantity - $1, updated_at = NOW() WHERE id = $2`,
		quantity, productID,
	)
	return err
}

func (r *OrderRepository) UpdateOrderTotal(ctx context.Context, orderID int, total decimal.Decimal) error {
	_, err := r.DB.ExecContext(
		ctx,
		`UPDATE orders SET total_amount = $1, status = $2, updated_at = NOW() WHERE id = $3`,
		total, 1, orderID,
	)
	return err
}
