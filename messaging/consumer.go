package messaging

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewConsumer() (*Consumer, error) {
	username := os.Getenv("RMQ_USERNAME")
	password := os.Getenv("RMQ_PASSWORD")
	host := os.Getenv("RMQ_HOST")
	port := 5672

	amqpURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		username,
		password,
		host,
		port,
	)

	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"price_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (rmq *Consumer) PriceUpdateConsume(ctx context.Context, handler func(ctx context.Context, d amqp091.Delivery) error) error {
	msgs, err := rmq.ch.Consume(
		"price_updates",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
		return err
	}
	go func() {
		log.Printf("Consuming....")
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			err := handler(ctx, msg)
			if err != nil {
				log.Fatalf("Failed to consume messages: will requeque %v", err)
				msg.Nack(false, true)
			}
		}
	}()

	return nil
}
