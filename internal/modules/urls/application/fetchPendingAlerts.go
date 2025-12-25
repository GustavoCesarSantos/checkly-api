package application

import (
	"context"
	"log/slog"

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
		slog.Error("Failed to fetch pending alerts", "error", err)
		return nil, err
	}
	if len(alerts) == 0 {
		slog.Info("No pending alerts found")
		return nil, utils.ErrRecordNotFound
	}
	slog.Info("Pending alerts found", "count", len(alerts))
	return alerts, nil
}