package database

import (
	"database/sql"
	"fmt"
	"orders-service/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	dsn := cfg.GetPostgreSQLDSN()
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			order_id VARCHAR(255) UNIQUE NOT NULL,
			user_id BIGINT NOT NULL,
			status VARCHAR(50) NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id BIGINT NOT NULL,
			quantity INTEGER NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_order_id ON orders(order_id)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetConn() *sql.DB {
	return db.conn
}
