# Orders Service

Сервис для обработки заказов, который потребляет сообщения из RabbitMQ, сохраняет их в PostgreSQL и кеширует в Redis.

## Архитектура

- **RabbitMQ**: Очередь сообщений с заказами
- **PostgreSQL**: Основное хранилище заказов
- **Redis**: Кеш для быстрого доступа к заказам

## Структура проекта

```
orders-service/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── models/              # Модели данных
│   │   └── order.go
│   ├── config/              # Конфигурация
│   │   └── config.go
│   ├── database/            # Подключение к PostgreSQL
│   │   └── postgres.go
│   ├── repository/          # Репозиторий для работы с БД
│   │   └── order_repository.go
│   ├── cache/               # Кеш Redis
│   │   └── redis_cache.go
│   ├── rabbitmq/            # Консьюмер RabbitMQ
│   │   └── consumer.go
│   └── service/             # Бизнес-логика
│       └── order_service.go
├── go.mod
└── README.md
```

## Требования

- Go 1.25.2 или выше
- PostgreSQL
- Redis
- RabbitMQ

## Установка и запуск

### 1. Клонирование и установка зависимостей

```bash
go mod download
```

### 2. Настройка переменных окружения

Создайте файл `.env` или установите переменные окружения:

```bash
# PostgreSQL
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=orders_db
export POSTGRES_SSLMODE=disable

# Redis
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=

# RabbitMQ
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USER=guest
export RABBITMQ_PASSWORD=guest
export RABBITMQ_QUEUE=orders

# Server
export SERVER_PORT=8080
```

### 3. Запуск с Docker Compose

Для удобства можно использовать `docker-compose.yml` для запуска всех зависимостей:

```bash
docker-compose up -d
```

### 4. Запуск сервиса

```bash
go run cmd/main.go
```

## Формат сообщений в RabbitMQ

Сервис ожидает JSON сообщения следующего формата:

```json
{
  "order_id": "ORD-12345",
  "user_id": 1001,
  "status": "pending",
  "total_price": 99.99,
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 49.99,
      "name": "Product Name"
    }
  ]
}
```

## Пример отправки сообщения в RabbitMQ

```bash
# Используя rabbitmqadmin
rabbitmqadmin publish exchange=amq.default routing_key=orders payload='{"order_id":"ORD-12345","user_id":1001,"status":"pending","total_price":99.99,"items":[{"product_id":1,"quantity":2,"price":49.99,"name":"Product Name"}]}'
```

## API

Сервис обрабатывает заказы автоматически из очереди RabbitMQ. При получении сообщения:

1. Заказ сохраняется в PostgreSQL
2. Заказ кешируется в Redis (TTL: 24 часа)

## База данных

Сервис автоматически создает необходимые таблицы при первом запуске:

- `orders` - основная таблица заказов
- `order_items` - товары в заказе

## Кеширование

Заказы кешируются в Redis с ключом `order:{order_id}` и TTL 24 часа.

