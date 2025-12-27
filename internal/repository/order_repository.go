package repository

import (
	"database/sql"
	"fmt"
	"orders-service/internal/models"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (order_id, user_id, status, total_price, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = tx.QueryRow(
		query,
		order.OrderID,
		order.UserID,
		order.Status,
		order.TotalPrice,
		time.Now(),
		time.Now(),
	).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	for i := range order.Items {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price, name)
					  VALUES ($1, $2, $3, $4, $5) RETURNING id`

		err = tx.QueryRow(
			itemQuery,
			order.ID,
			order.Items[i].ProductID,
			order.Items[i].Quantity,
			order.Items[i].Price,
			order.Items[i].Name,
		).Scan(&order.Items[i].ID)

		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
		order.Items[i].OrderID = order.ID
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetOrderByOrderID(orderID string) (*models.Order, error) {
	order := &models.Order{}

	query := `SELECT id, order_id, user_id, status, total_price, created_at, updated_at
			  FROM orders WHERE order_id = $1`

	err := r.db.QueryRow(query, orderID).Scan(
		&order.ID,
		&order.OrderID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	itemsQuery := `SELECT id, order_id, product_id, quantity, price, name
				   FROM order_items WHERE order_id = $1`

	rows, err := r.db.Query(itemsQuery, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Name,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderRepository) UpdateOrderStatus(orderID string, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE order_id = $3`
	_, err := r.db.Exec(query, status, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
