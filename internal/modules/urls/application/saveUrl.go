package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

type SaveUrl struct {
	urlRepository db.IUrlRepository
}

func NewSaveUrl(urlRepository db.IUrlRepository) *SaveUrl {
	return &SaveUrl{
		urlRepository,
	}
}

func (su *SaveUrl) Execute(input dtos.CreateUrlRequest) (*domain.Url, error) {
	url := domain.NewUrl(
		input.Url, 
		input.Interval, 
		input.RetryLimit, 
		input.Contact,
		*input.Status,
		*input.NextCheck,
	)
	err := su.urlRepository.Save(url)
	if(err != nil) {
		return nil, err
	}
	return url, nil
}