package application

import (
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"context"
)

type UpdateUrl struct {
	urlRepository db.IUrlRepository
}

func NewUpdateUrl(urlRepository db.IUrlRepository) *UpdateUrl {
	return &UpdateUrl{
		urlRepository,
	}
}

func (u *UpdateUrl) Execute(ctx context.Context, urlId int64, input dtos.UpdateUrlRequest) error {
	params := db.UpdateUrlParams{
		NextCheck:      input.NextCheck,
		RetryCount:     input.RetryCount,
		StabilityCount: input.StabilityCount,
		Status:         input.Status,
	}
	err := u.urlRepository.Update(ctx, urlId, params)
	if err != nil {
		return err
	}
	return nil
}
