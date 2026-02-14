package application

import (
	"context"
	"fmt"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type FetchPendingAlerts struct {
	alertOutboxRepositroy db.IAlertOutboxRepository
}

func NewFetchPendingAlerts(alertOutboxRepositroy db.IAlertOutboxRepository) *FetchPendingAlerts {
	return &FetchPendingAlerts{
		alertOutboxRepositroy: alertOutboxRepositroy,
	}
}

func (f *FetchPendingAlerts) Execute(ctx context.Context, limit int) ([]domain.AlertOutbox, error) {
	alerts, err := f.alertOutboxRepositroy.FindAllPendingAlerts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("fetchPendingAlerts: %w", err)
	}
	if len(alerts) == 0 {
		return nil, fmt.Errorf("fetchPendingAlerts: %w", utils.ErrRecordNotFound)
	}
	return alerts, nil
}
