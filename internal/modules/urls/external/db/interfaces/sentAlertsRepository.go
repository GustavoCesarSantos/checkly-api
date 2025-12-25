package db

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"context"
)

type UpdateAlertParams struct {
	Status *domain.AlertStatus
}

type ISentAlertsRepository interface {
	Save(idempotencyKey string) error
	Update(ctx context.Context, idempotencyKey string, params UpdateAlertParams) error
}
