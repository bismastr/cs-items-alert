package item

import (
	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

type ItemService struct {
	timescaleRepo timescale_repository.Repository
	postgresRepo  repository.Repository
}

func NewItemService(timescaleRepo timescale_repository.Repository, postgresRepo repository.Repository) *ItemService {
	return &ItemService{
		timescaleRepo: timescaleRepo,
		postgresRepo:  postgresRepo,
	}
}
