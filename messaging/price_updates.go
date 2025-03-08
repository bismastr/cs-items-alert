package messaaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (p *Publisher) PublishPriceUpdate(itemId int) error {
	body := fmt.Sprintf(`{"item_id": %d}`, itemId)

	return p.ch.Publish(
		"",
		"price_updates",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
}
