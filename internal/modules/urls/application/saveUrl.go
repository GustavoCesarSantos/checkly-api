package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"time"
)

type SaveUrl struct {
	urlRepository db.IUrlRepository
}

func NewSaveUrl(urlRepository db.IUrlRepository) *SaveUrl {
	return &SaveUrl{
		urlRepository,
	}
}

func (su *SaveUrl) Execute(input dtos.CreateUrlRequest, isHealthy bool) (*domain.Url, error) {
	status := domain.StatusHealthy
	nextCheck := time.Now().Add(time.Duration(input.Interval) * time.Minute)
	if(!isHealthy) {
		status = domain.StatusDegraded
		nextCheck = time.Now().Add(time.Minute)
	}
	url := domain.NewUrl(
		input.Url, 
		input.Interval, 
		input.RetryLimit, 
		input.Contact,
		status,
		nextCheck,
	)
	err := su.urlRepository.Save(url)
	if(err != nil) {
		return nil, err
	}
	return url, nil
}