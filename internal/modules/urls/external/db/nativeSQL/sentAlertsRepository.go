package nativeSQL

import (
	"context"
	"database/sql"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type sentAlertsRepository struct {
	DB *sql.DB
}

func NewSentAlertsRepository(db *sql.DB) db.ISentAlertsRepository {
	return &sentAlertsRepository{
		DB: db,
	}
}

func (a *sentAlertsRepository) Save(ctx context.Context, idempotencyKey string) error {
	query := `
        INSERT INTO sent_alerts (
            idempotency_key
        )
        VALUES (
            $1
        )
    `
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return a.DB.QueryRowContext(queryCtx, query, idempotencyKey).Err()
}

func (a *sentAlertsRepository) Update(ctx context.Context, idempotencyKey string, status domain.AlertStatus) error {
	query := `
		UPDATE 
			sent_alerts 
		SET 
			status = $1, 
			sent_at = NOW() 
		WHERE 
			idempotency_key = $2
	`
	args := []any{
		status,
		idempotencyKey,
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
