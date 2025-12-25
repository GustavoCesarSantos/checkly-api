package db

import (
	"context"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type UpdateAlertParams struct {
	Status *domain.AlertStatus
}

type ISentAlertsRepository interface {
	Save(idempotencyKey string) error
	Update(ctx context.Context, idempotencyKey string, status domain.AlertStatus) error
}
