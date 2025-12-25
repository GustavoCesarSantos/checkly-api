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

func (a *sentAlertsRepository) Save(idempotencyKey string) error {
	query := `
        INSERT INTO sent_alerts (
            idempotency_key
        )
        VALUES (
            $1
        )
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return a.DB.QueryRowContext(ctx, query, idempotencyKey).Err()
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
	result, err := a.DB.ExecContext(ctx, query, args...)
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
