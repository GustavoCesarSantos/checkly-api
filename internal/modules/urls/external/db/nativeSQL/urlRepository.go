package db

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"context"
	"database/sql"
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
            interval,
            retry_limit,
            contact
        )
        VALUES (
            $1,
            $2,
            $3
        )
        RETURNING
            id,
			created_at
    `
    args := []any{
        url.Interval,
		url.RetryLimit,
		url.Contact,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
    return ur.DB.QueryRowContext(ctx, query, args...).Scan(
        &url.ID,
        &url.CreatedAt,
    )
}