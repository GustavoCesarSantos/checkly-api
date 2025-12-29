package application

import (
	"context"
	"fmt"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	factory "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/factory/repositoryFactory/interface"
	unitOfWork "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/unitOfWork/repositoryUnitOfWork/interface"
)

type UpdateUrlWithOutbox struct {
	repositoryUnitOfWork unitOfWork.IRepositoryUnitOfWork
}

func NewUpdateUrlWithOutbox(repositoryUnitOfWork unitOfWork.IRepositoryUnitOfWork) *UpdateUrlWithOutbox {
	return &UpdateUrlWithOutbox{
		repositoryUnitOfWork: repositoryUnitOfWork,
	}
}

func (u *UpdateUrlWithOutbox) Execute(ctx context.Context, url domain.Url, input dtos.UpdateUrlRequest) error {
	return u.repositoryUnitOfWork.WithTx(ctx, func(f factory.IRepositoryFactory) error {
		params := db.UpdateUrlParams{
			NextCheck:      input.NextCheck,
			RetryCount:     input.RetryCount,
			DownCount:      input.DownCount,
			StabilityCount: input.StabilityCount,
			Status:         input.Status,
		}
		updateErr := f.Urls().Update(ctx, url.ID, params)
		if updateErr != nil {
			return fmt.Errorf("updateUrlWithOutbox: %w", updateErr)
		}
		alert := domain.NewAlertOutbox(
			url.ID,
			domain.Payload{
				Url:   url.Address,
				Email: url.Contact,
			},
		)
		saveErr := f.AlertOutbox().Save(ctx, alert)
		if saveErr != nil {
			return fmt.Errorf("updateUrlWithOutbox: %w", saveErr)
		}
		return nil
	})
}