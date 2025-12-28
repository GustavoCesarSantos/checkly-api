package application

import (
	"context"
	"fmt"
	"time"

	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
)

type MarkSent struct {
	alertOutboxRepo db.IAlertOutboxRepository
}

func NewMarkSent(alertOutboxRepo db.IAlertOutboxRepository) *MarkSent {
	return &MarkSent{
		alertOutboxRepo: alertOutboxRepo,
	}
}

func (m *MarkSent) Execute(ctx context.Context, alertID int64, sentAt time.Time) error {
	err := m.alertOutboxRepo.Update(ctx, alertID, sentAt)
	if err != nil {
		return fmt.Errorf("markSent: %w", err)
	}
	return nil
}