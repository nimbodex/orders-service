package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"orders-service/internal/config"
	"orders-service/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	url := cfg.GetRabbitMQURL()

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	queue := cfg.RabbitMQ.Queue

	_, err = channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (c *Consumer) Consume(handler func(*models.Order) error) error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack (false - ручное подтверждение)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			var order models.Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				log.Printf("Error unmarshaling order: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := handler(&order); err != nil {
				log.Printf("Error processing order: %v", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	return nil
}

func (c *Consumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
