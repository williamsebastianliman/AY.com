package service

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
)

type EventBus interface {
  Publish(ctx context.Context, topic string, exchange string, payload interface{}) error
}

type RabbitEventBus struct {
  conn *amqp.Connection
}

func NewRabbitEventBus(url string) (*RabbitEventBus, error) {
  conn, err := amqp.Dial(url)
  if err != nil {
    return nil, err
  }
  return &RabbitEventBus{conn: conn}, nil
}

func (b *RabbitEventBus) Publish(ctx context.Context, topic string, exchange string, payload interface{}) error {
  ch, err := b.conn.Channel()
  if err != nil {
    return err
  }
  defer ch.Close()

  body, err := json.Marshal(payload)
  if err != nil {
    return err
  }
  return ch.Publish(
    exchange,    
    topic, 
    false, 
    false,
    amqp.Publishing{
      ContentType: "application/json",
      Body:        body,
    },
  )
}