package application

import (
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
	"context"
	"errors"
	"log/slog"
	"time"
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
		switch {
		case errors.Is(err, utils.ErrRecordNotFound):
			slog.Warn("alert not found for mark as sent", "alertID", alertID)
		default:
			slog.Error("failed to mark alert as sent", "alertID", alertID, "error", err)
		}
		return err
	}
	return nil
}