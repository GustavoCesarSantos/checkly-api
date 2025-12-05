package db

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"time"
)

type UpdateUrlParams struct {
	NextCheck *time.Time
	RetryCount *int
}

type IUrlRepository interface {
	Save(url *domain.Url) error
	Update(urlId int64, params UpdateUrlParams) error
}