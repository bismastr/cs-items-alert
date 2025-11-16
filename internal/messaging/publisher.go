package messaging

import (
	"context"
	"sync"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn           *amqp091.Connection
	channel        *amqp091.Channel
	declaredQueues map[string]bool
	mu             sync.Mutex
}

func NewPublisher(config *config.Config) (*Publisher, error) {
	conn, err := amqp091.Dial(config.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Publisher{
		conn:           conn,
		channel:        ch,
		declaredQueues: make(map[string]bool),
	}, nil
}

func (p *Publisher) EnsureQueue(q string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.declaredQueues[q] {
		return nil
	}

	_, err := p.channel.QueueDeclare(
		q,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	p.declaredQueues[q] = true
	return nil
}

func (p *Publisher) Publish(ctx context.Context, queueName string, message []byte) error {
	err := p.EnsureQueue(queueName)
	if err != nil {
		return err
	}

	payload := amqp091.Publishing{
		ContentType: "application/json",
		Body:        message,
	}

	err = p.channel.PublishWithContext(
		ctx,
		"",        //exchange (default)
		queueName, //routing
		false,     //mandotaory
		false,     //immediate
		payload)
	if err != nil {
		return err
	}

	return nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}
