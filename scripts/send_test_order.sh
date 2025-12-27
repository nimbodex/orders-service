#!/bin/bash

RABBITMQ_HOST=${RABBITMQ_HOST:-localhost}
RABBITMQ_PORT=${RABBITMQ_PORT:-5672}
RABBITMQ_USER=${RABBITMQ_USER:-guest}
RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD:-guest}
RABBITMQ_QUEUE=${RABBITMQ_QUEUE:-orders}

ORDER_JSON='{
  "order_id": "ORD-'$(date +%s)'",
  "user_id": 1001,
  "status": "pending",
  "total_price": 199.99,
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 49.99,
      "name": "Product A"
    },
    {
      "product_id": 2,
      "quantity": 1,
      "price": 100.01,
      "name": "Product B"
    }
  ]
}'

echo "Sending test order to RabbitMQ..."
echo "$ORDER_JSON" | docker run --rm -i --network host \
  rabbitmq:3-management-alpine \
  rabbitmqadmin -H $RABBITMQ_HOST -P 15672 -u $RABBITMQ_USER -p $RABBITMQ_PASSWORD \
  publish exchange=amq.default routing_key=$RABBITMQ_QUEUE payload="$ORDER_JSON"

echo "Order sent successfully!"

