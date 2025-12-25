package nativeSQL

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	urls_utils "GustavoCesarSantos/checkly-api/internal/modules/urls/utils"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
	"context"
	"time"
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
			*
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
	rows, err := a.DB.QueryContext(ctx, query, limit)
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
			&alert.IdempotencyKey,
			&alert.SentAt,
			&alert.ProcessingAt,
			&alert.RetryCount,
			&alert.NextRetryAt,
			&alert.LockedAt,
			&alert.LockedBy,
			&alert.CreatedAt,
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

func (a *alertOutboxRepository) Save(alert *domain.AlertOutbox) error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return a.DB.QueryRowContext(ctx, query, args...).Err()
}

func (a *alertOutboxRepository) UpdateRetryInfo(ctx context.Context, id int64) error {
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
	result, err := a.DB.ExecContext(ctx, query, id)
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