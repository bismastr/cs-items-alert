package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetAlertsRealtime struct {
	Id            int
	DiscordId     int64
	ItemId        int
	ConditionType string
	IsActive      bool
	Threshold     float64
}

const getAlertsRealTime = `
SELECT 
	id,
	discord_id,
	item_id,
	condition_type,
	is_active,
	threshold 
	FROM alerts_real_time
WHERE is_active = true
`

func (q *Queries) GetAlertsRealtime(ctx context.Context) ([]GetAlertsRealtime, error) {
	rows, err := q.db.Query(ctx, getAlertsRealTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GetAlertsRealtime
	for rows.Next() {
		var i GetAlertsRealtime
		rows.Scan(
			&i.Id,
			&i.DiscordId,
			&i.ItemId,
			&i.ConditionType,
			&i.IsActive,
			&i.Threshold,
		)

		items = append(items, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

type GetDailyAlertSchedule struct {
	Id        int
	ItemId    int
	IsActive  bool
	CreatedAt pgtype.Timestamp
	DiscordId int64
}

const getDailyAlertSchedule = `
SELECT * FROM alerts_daily_schedule
WHERE is_active = true
`

func (q *Queries) GetDailyAlertSchedule(ctx context.Context) ([]GetDailyAlertSchedule, error) {
	rows, err := q.db.Query(ctx, getDailyAlertSchedule)
	if err != nil {
		return nil, err
	}

	var items []GetDailyAlertSchedule
	for rows.Next() {
		var i GetDailyAlertSchedule
		rows.Scan(
			&i.Id,
			&i.ItemId,
			&i.IsActive,
			&i.CreatedAt,
			&i.DiscordId,
		)

		err := rows.Err()
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}
