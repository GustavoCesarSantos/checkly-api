package application

import (
	"context"
	"fmt"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type FetchUrls struct {
	urlRepository db.IUrlRepository
}

func NewFetchUrls(urlRepository db.IUrlRepository) *FetchUrls {
	return &FetchUrls{
		urlRepository: urlRepository,
	}
}

func (f *FetchUrls) Execute(ctx context.Context, nextCheck time.Time) ([]domain.Url, error) {
	urls, err := f.urlRepository.FindAllByNextCheck(ctx, nextCheck)
	if err != nil {
		return nil, fmt.Errorf("fetchUrls: %w", err)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("fetchUrls: %w", utils.ErrRecordNotFound)
	}
	return urls, nil
}