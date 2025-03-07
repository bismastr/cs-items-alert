package messaaging

import (
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type RmqClient struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRmqClient() (*RmqClient, error) {
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

	return &RmqClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *RmqClient) Close() {
	p.ch.Close()
	p.conn.Close()
}
