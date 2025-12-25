package nativeSQL

import (
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
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

func (a *sentAlertsRepository) Update(ctx context.Context, idempotencyKey string, params db.UpdateAlertParams) error {
	if params.Status == nil {
		return errors.New("NO COLUMN FIELD PROVIDED FOR UPDATING")
	}
	query := "UPDATE sent_alerts SET"
	var args []interface{}
	argPos := 1
	if params.Status != nil {
		query += " status = $" + strconv.Itoa(argPos) + ","
		args = append(args, *params.Status)
		argPos++
	}
	query += " sent_at = NOW()"
	query = strings.TrimSuffix(query, ",") + " WHERE idempotency_key = $" + strconv.Itoa(argPos)
	args = append(args, idempotencyKey)
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
