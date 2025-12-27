package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"orders-service/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO orders (order_id, status, amount)
		VALUES ($1, $2, $3)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		order.OrderID,
		order.Status,
		order.TotalPrice,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrderByOrderID(ctx context.Context, orderID string) (*models.Order, error) {
	query := `
		SELECT order_id, status, amount
		FROM orders
		WHERE order_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, orderID)

	var order models.Order
	if err := row.Scan(&order.OrderID, &order.Status, &order.TotalPrice); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE order_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, orderID)
	return err
}
