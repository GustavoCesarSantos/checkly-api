package application

import (
	"context"
	"errors"
	"log/slog"

	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type ScheduleRetryAlert struct {
	AlertOutboxRepo db.IAlertOutboxRepository
}

func NewScheduleRetryAlert(alertOutboxRepo db.IAlertOutboxRepository) *ScheduleRetryAlert {
	return &ScheduleRetryAlert{
		AlertOutboxRepo: alertOutboxRepo,
	}
}

func (s *ScheduleRetryAlert) Execute(ctx context.Context, alertId int64) error {
	err := s.AlertOutboxRepo.UpdateRetryInfo(ctx, alertId)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrRecordNotFound):
			slog.Warn("alert not found for retry scheduling", "alertID", alertId)
		default:
			slog.Error("failed to schedule retry for alert", "alertID", alertId, "error", err)
		}
		return err
	}
	return nil
}