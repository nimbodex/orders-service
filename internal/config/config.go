package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgreSQL PostgreSQLConfig
	Redis      RedisConfig
	RabbitMQ   RabbitMQConfig
	Server     ServerConfig
}

type PostgreSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Queue    string
}

type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	config := &Config{
		PostgreSQL: PostgreSQLConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "orders_db"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			Queue:    getEnv("RABBITMQ_QUEUE", "orders"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}

	return config, nil
}

func (c *Config) GetPostgreSQLDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgreSQL.Host,
		c.PostgreSQL.Port,
		c.PostgreSQL.User,
		c.PostgreSQL.Password,
		c.PostgreSQL.DBName,
		c.PostgreSQL.SSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) GetRabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
