package db

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"context"
	"time"
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
}
