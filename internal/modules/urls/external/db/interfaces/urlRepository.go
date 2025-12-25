package db

import (
	"context"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type UpdateUrlParams struct {
	NextCheck      *time.Time
	RetryCount     *int
	StabilityCount *int
	Status         *domain.UrlStatus
}

type IUrlRepository interface {
	FindAllByNextCheck(ctx context.Context, nextCheck time.Time) ([]domain.Url, error)
	Save(url *domain.Url) error
	Update(ctx context.Context, urlId int64, params UpdateUrlParams) error
	UpdateToNotified(ctx context.Context, urlId int64) error
}
