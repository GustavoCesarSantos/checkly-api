package nativeSQL

import (
	"context"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	urls_utils "GustavoCesarSantos/checkly-api/internal/modules/urls/utils"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type alertOutboxRepository struct {
	DB urls_utils.DBExecutor
}

func NewAlertOutboxRepository(db urls_utils.DBExecutor) db.IAlertOutboxRepository {
	return &alertOutboxRepository{
		DB: db,
	}
}

func (a *alertOutboxRepository) FindAllPendingAlerts(ctx context.Context, limit int) ([]domain.AlertOutbox, error) {
	query := `
        SELECT 
			id,
			url_id,
			payload
		FROM 
			alert_outbox
		WHERE 
			sent_at IS NULL
			AND next_retry_at <= now()
			AND (locked_at IS NULL OR locked_at < now() - interval '5 minutes')
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
    `
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := a.DB.QueryContext(queryCtx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	alerts := []domain.AlertOutbox{}
	for rows.Next() {
		var alert domain.AlertOutbox
		err := rows.Scan(
			&alert.ID,
			&alert.UrlId,
			&alert.Payload,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	return alerts, nil
}

func (a *alertOutboxRepository) Save(ctx context.Context, alert *domain.AlertOutbox) error {
	query := `
        INSERT INTO alert_outbox (
            url_id,
			payload
        )
        VALUES (
            $1,
            $2
        )
    `
	args := []any{
		alert.UrlId,
		alert.Payload,
	}
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return a.DB.QueryRowContext(queryCtx, query, args...).Err()
}

func (a *alertOutboxRepository) Update(ctx context.Context, alertId int64, sentAt time.Time) error {
	query := `
		UPDATE 
			alert_outbox
		SET 
			sent_at = $1,
			updated_at = NOW()
		WHERE 
			id = $2
	`
	args := []any{
		sentAt,
		alertId,
	}
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	result, err := a.DB.ExecContext(queryCtx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, rowsAffectedErr := result.RowsAffected()
	if rowsAffectedErr != nil {
		return rowsAffectedErr
	}
	if rowsAffected == 0 {
		return utils.ErrRecordNotFound
	}
	return nil
}

func (a *alertOutboxRepository) UpdateRetryInfo(ctx context.Context, alertId int64) error {
	query := `
		UPDATE 
			alert_outbox
        SET 
			retry_count = retry_count + 1,
            next_retry_at = now() + (interval '30 seconds' * power(2, retry_count)),
			updated_at = NOW()
        WHERE 
			id = $1
	`
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	result, err := a.DB.ExecContext(queryCtx, query, alertId)
	if err != nil {
		return err
	}
	rowsAffected, rowsAffectedErr := result.RowsAffected()
	if rowsAffectedErr != nil {
		return rowsAffectedErr
	}
	if rowsAffected == 0 {
		return utils.ErrRecordNotFound
	}
	return nil
}