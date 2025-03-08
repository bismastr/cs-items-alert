package messaaging

import (
	"encoding/json"
	"log"
)

func (rmq *RmqClient) PriceUpdateConsume() {
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
	}

	for msg := range msgs {
		var priceUpdate struct {
			ItemId int `json:"item_id"`
		}
		if err := json.Unmarshal(msg.Body, &priceUpdate); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}
	}
}
