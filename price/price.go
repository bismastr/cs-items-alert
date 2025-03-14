package price

import (
	"context"
	"fmt"
	"log"

	messaaging "github.com/bismastr/cs-price-alert/messaging"
	"github.com/bismastr/cs-price-alert/repository"
)

type PriceService struct {
	repo      *repository.Queries
	publisher *messaaging.Publisher
}

func NewPriceService(repo *repository.Queries, publihser *messaaging.Publisher) *PriceService {
	return &PriceService{repo: repo, publisher: publihser}
}

func (s *PriceService) InsertItem(ctx context.Context, item repository.InsertItem) error {
	id, err := s.repo.InsertItem(ctx, item)
	if err != nil {
		return err
	}

	message := fmt.Sprintf(`{"item_id": %d}`, id)

	err = s.publisher.Publish("price_updates", []byte(message))
	if err != nil {
		log.Printf("Failed to publish price update: %v", err)
		return err
	}

	return nil
}
