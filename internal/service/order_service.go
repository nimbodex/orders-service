package service

import (
	"context"
	"fmt"
	"log"
	"orders-service/internal/cache"
	"orders-service/internal/models"
	"orders-service/internal/repository"
	"time"
)

type OrderService struct {
	repo  *repository.OrderRepository
	cache *cache.Cache
}

func NewOrderService(repo *repository.OrderRepository, cache *cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	log.Printf("Order saved to database: %s (ID: %d)", order.OrderID, order.ID)

	if err := s.cache.SetOrder(ctx, order, 24*time.Hour); err != nil {
		log.Printf("Warning: failed to cache order: %v", err)
	} else {
		log.Printf("Order cached in Redis: %s", order.OrderID)
	}

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	order, err := s.cache.GetOrder(ctx, orderID)
	if err == nil {
		log.Printf("Order retrieved from cache: %s", orderID)
		return order, nil
	}

	order, err = s.repo.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := s.cache.SetOrder(ctx, order, 24*time.Hour); err != nil {
		log.Printf("Warning: failed to cache order after retrieval: %v", err)
	}

	log.Printf("Order retrieved from database: %s", orderID)
	return order, nil
}
