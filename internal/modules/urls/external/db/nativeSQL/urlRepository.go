package nativeSQL

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	urls_utils "GustavoCesarSantos/checkly-api/internal/modules/urls/utils"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type urlRepository struct {
	DB urls_utils.DBExecutor
}

func NewUrlRepository(db urls_utils.DBExecutor) db.IUrlRepository {
	return &urlRepository{
		DB: db,
	}
}

func (u *urlRepository) FindAllByNextCheck(ctx context.Context, nextCheck time.Time) ([]domain.Url, error) {
	query := `
        SELECT
            id,
            address,
            interval,
            retry_limit,
            retry_count,
            stability_count,
            contact,
            next_check,
            status
        FROM
            urls
        WHERE
            next_check <= $1
            AND status NOT IN (30, 40);
    `
	rows, err := u.DB.QueryContext(ctx, query, nextCheck)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := []domain.Url{}
	for rows.Next() {
		var url domain.Url
		err := rows.Scan(
			&url.ID,
			&url.Address,
			&url.Interval,
			&url.RetryLimit,
			&url.RetryCount,
			&url.StabilityCount,
			&url.Contact,
			&url.NextCheck,
			&url.Status,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	return urls, nil
}

func (u *urlRepository) Save(url *domain.Url) error {
	query := `
        INSERT INTO urls (
            address,
            interval,
            retry_limit,
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
            $6
        )
        RETURNING
            external_id,
			created_at
    `
	args := []any{
		url.Address,
		url.Interval,
		url.RetryLimit,
		url.Contact,
		url.Status,
		url.NextCheck,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return u.DB.QueryRowContext(ctx, query, args...).Scan(
		&url.ExternalID,
		&url.CreatedAt,
	)
}

func (u *urlRepository) Update(ctx context.Context, urlId int64, params db.UpdateUrlParams) error {
	if params.NextCheck == nil &&
		params.RetryCount == nil &&
		params.StabilityCount == nil &&
		params.Status == nil {
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
	if params.StabilityCount != nil {
		query += " stability_count = $" + strconv.Itoa(argPos) + ","
		args = append(args, *params.StabilityCount)
		argPos++
	}
	if params.Status != nil {
		query += " status = $" + strconv.Itoa(argPos) + ","
		args = append(args, *params.Status)
		argPos++
	}
	query += " updated_at = NOW()"
	query = strings.TrimSuffix(query, ",") + " WHERE id = $" + strconv.Itoa(argPos)
	args = append(args, urlId)
	result, err := u.DB.ExecContext(ctx, query, args...)
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

func (u *urlRepository) UpdateToNotified(ctx context.Context, urlId int64) error {
	query := `
		UPDATE 
			urls
		SET 
			status = 40,
			updated_at = NOW()
		WHERE 
			id = $1
			AND status = 30;
	`
	result, err := u.DB.ExecContext(ctx, query, urlId)
	if err != nil {
		return err
	}
	rowsAffected, rowsAffectedErr := result.RowsAffected()
	if rowsAffectedErr != nil {
		return rowsAffectedErr
	}
	if rowsAffected == 0 {
		return nil
	}
	return nil
}