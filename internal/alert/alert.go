package alert

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

// TODO: i want to implement daily alert if price above certain threshold
type TimescaleRepository interface {
	Get24HourPricesChanges(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error)
}

type AlertService struct {
	timescaleRepository TimescaleRepository
}

func NewAlertService(timescaleRepo TimescaleRepository) *AlertService {
	return &AlertService{
		timescaleRepository: timescaleRepo,
	}
}

func (s *AlertService) Alert24Hour(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error) {
	return s.timescaleRepository.Get24HourPricesChanges(ctx)
}
