package messaging

import (
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn *amqp091.Connection
}

func NewPublihser(url string) (*Publisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn: conn,
	}, nil
}

func (p *Publisher) Publish(q string, message []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	ch.QueueDeclare(q, true, false, false, false, nil)

	payload := amqp091.Publishing{
		ContentType: "application/json",
		Body:        message,
	}

	err = ch.Publish("", q, false, false, payload)
	if err != nil {
		return err
	}

	return nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}
