package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"orders-service/internal/config"
	"orders-service/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: client}, nil
}

func (c *Cache) SetOrder(ctx context.Context, order *models.Order, ttl time.Duration) error {
	key := fmt.Sprintf("order:%s", order.OrderID)

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set order in cache: %w", err)
	}

	return nil
}

func (c *Cache) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	key := fmt.Sprintf("order:%s", orderID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("order not found in cache")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order from cache: %w", err)
	}

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}

func (c *Cache) DeleteOrder(ctx context.Context, orderID string) error {
	key := fmt.Sprintf("order:%s", orderID)
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete order from cache: %w", err)
	}
	return nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}
