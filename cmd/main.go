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
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

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

	err = consumer.Consume(func(order *models.Order) error {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		return orderService.ProcessOrder(ctxWithTimeout, order)
	})

	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	<-ctx.Done()
	log.Println("Shutting down server...")
}
