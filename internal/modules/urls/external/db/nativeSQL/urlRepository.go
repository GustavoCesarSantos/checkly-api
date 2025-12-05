package db

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

type urlRepository struct {
	DB *sql.DB
}

func NewUrlRepository(db *sql.DB) db.IUrlRepository {
	return &urlRepository {
		DB: db,
	}
}

func (ur *urlRepository) Save(url *domain.Url) error {
	query := `
        INSERT INTO urls (
            url,
            interval,
            retry_limit,
            retry_count,
            contact,
            status,
            next_check
        )
        VALUES (
            $1,
            $2,
            $3,
            $4,
            $5,
            $6,
            $7
        )
        RETURNING
            id,
			created_at
    `
    args := []any{
        url.Url,
        url.Interval,
		url.RetryLimit,
		url.RetryCount,
        url.Contact,
        url.Status,
        url.NextCheck,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
    return ur.DB.QueryRowContext(ctx, query, args...).Scan(
        &url.ID,
        &url.CreatedAt,
    )
}

func (ur *urlRepository) Update(urlId int64, params db.UpdateUrlParams) error {
    if params.NextCheck == nil && params.RetryCount == nil {
		return errors.New("NO COLUMN FIELD PROVIDED FOR UPDATING")
	}
    query := "UPDATE urls SET"
	var args []interface{}
	argPos := 1
    if params.NextCheck != nil {
		query += " next_check = $" + strconv.Itoa(argPos) + ","
		args = append(args, *params.NextCheck)
		argPos++
	}
	if params.RetryCount != nil {
		query += " retry_count = $" + strconv.Itoa(argPos) + ","
		args = append(args, *params.RetryCount)
		argPos++
	}
    query += " updated_at = NOW()"
    query = strings.TrimSuffix(query, ",") + " WHERE id = $" + strconv.Itoa(argPos)
	args = append(args, urlId)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := ur.DB.ExecContext(ctx, query, args...)
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