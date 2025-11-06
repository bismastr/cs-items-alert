package analysis

import (
	"context"
	"github.com/bismastr/cs-price-alert/internal/repository"
)

func (a *Analysis) alertsRealTime(ctx context.Context) (map[int]repository.GetAlertsRealtime, error) {
	alerts, err := a.repo.GetAlertsRealtime(ctx)
	if err != nil {
		return nil, err
	}

	alertsMap := make(map[int]repository.GetAlertsRealtime)
	for _, v := range alerts {
		alertsMap[v.ItemId] = v
	}

	return alertsMap, err
}
