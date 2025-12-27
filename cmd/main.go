package main

import (
	"context"
	"log"
	"orders-service/internal/cache"
	"orders-service/internal/config"
	"orders-service/internal/database"
	"orders-service/internal/models"
	"orders-service/internal/rabbitmq"
	"orders-service/internal/repository"
	"orders-service/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	redisCache, err := cache.NewCache(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()
	log.Println("Connected to Redis")

	orderRepo := repository.NewOrderRepository(db.GetConn())
	orderService := service.NewOrderService(orderRepo, redisCache)

	consumer, err := rabbitmq.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer consumer.Close()
	log.Println("Connected to RabbitMQ")

	if err := consumer.Consume(func(order *models.Order) error {
		ctx := context.Background()
		return orderService.ProcessOrder(ctx, order)
	}); err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
