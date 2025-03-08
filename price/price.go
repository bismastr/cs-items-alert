package price

import (
	"context"
	"log"

	messaaging "github.com/bismastr/cs-price-alert/messaging"
	"github.com/bismastr/cs-price-alert/repository"
)

type PriceService struct {
	repo      *repository.Queries
	publisher *messaaging.RmqClient
}

func NewPriceService(repo *repository.Queries, publihser *messaaging.RmqClient) *PriceService {
	return &PriceService{repo: repo, publisher: publihser}
}

func (s *PriceService) InsertItem(ctx context.Context, item repository.InsertItem) error {
	id, err := s.repo.InsertItem(ctx, item)
	if err != nil {
		return err
	}

	err = s.publisher.PublishPriceUpdate(id)
	if err != nil {
		log.Printf("Failed to publish price update: %v", err)
		return err
	}

	return nil
}
