package db

import (
	"context"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type IAlertOutboxRepository interface {
	FindAllPendingAlerts(ctx context.Context, limit int) ([]domain.AlertOutbox, error)
	Save(alert *domain.AlertOutbox) error
	UpdateRetryInfo(ctx context.Context, alertId int64) error
}
