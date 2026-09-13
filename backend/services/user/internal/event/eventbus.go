package event

import (
	"context"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type EventBus interface {
  Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, msg []byte)) error
  Publish(ctx context.Context, topic string, msg []byte) error
}

type RabbitMQEventBus struct {
  conn     *amqp.Connection
  channel  *amqp.Channel
  exchange string
}

func NewRabbitMQEventBus(conn *amqp.Connection) *RabbitMQEventBus {
  ch, err := conn.Channel()
  if err != nil {
    log.Fatalf("AMQP channel error: %v", err)
  }
  exchange := "user_events"
  if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
    log.Fatalf("ExchangeDeclare error: %v", err)
  }
  return &RabbitMQEventBus{conn: conn, channel: ch, exchange: exchange}
}

func (r *RabbitMQEventBus) Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, msg []byte)) error {
  q, err := r.channel.QueueDeclare("", false, true, true, false, nil)
  if err != nil {
    return err
  }
  if err := r.channel.QueueBind(q.Name, topic, r.exchange, false, nil); err != nil {
    return err
  }
  msgs, err := r.channel.Consume(q.Name, "", true, true, false, false, nil)
  if err != nil {
    return err
  }
  go func() {
    for d := range msgs {
      handler(ctx, d.Body)
    }
  }()
  return nil
}

func (r *RabbitMQEventBus) Publish(ctx context.Context, topic string, msg []byte) error {
  select {
  case <-ctx.Done():
    return ctx.Err()
  default:
  }
  return r.channel.Publish(r.exchange, topic, false, false, amqp.Publishing{
    Timestamp:   time.Now(),
    ContentType: "application/json",
    Body:        msg,
  })
}