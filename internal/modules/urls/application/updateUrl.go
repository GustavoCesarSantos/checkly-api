package application

import (
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

type UpdateUrl struct {
	urlRepository db.IUrlRepository
}

func NewUpdateUrl(urlRepository db.IUrlRepository) *UpdateUrl {
	return &UpdateUrl{
		urlRepository,
	}
}

func (uu *UpdateUrl) Execute(urlId int64, input dtos.UpdateUrlRequest) error {
	params := db.UpdateUrlParams{
		NextCheck: input.NextCheck,
	}
	err := uu.urlRepository.Update(urlId, params)
	if(err != nil) {
		return err
	}
	return nil
}