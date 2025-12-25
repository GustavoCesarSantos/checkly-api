package db

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"context"
)

type UpdateAlertParams struct {
	IdempotencyKey *string
}

type IAlertOutboxRepository interface {
	FindAllPendingAlerts(ctx context.Context, limit int) ([]domain.AlertOutbox, error)
	Save(alert *domain.AlertOutbox) error
	Update(ctx context.Context, alertId int64, params UpdateAlertParams) error
	UpdateRetryInfo(ctx context.Context, alertId int64) error
}
