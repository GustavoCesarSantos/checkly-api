package application

import (
	"context"
	"fmt"

	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
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
		return fmt.Errorf("scheduleRetryAlert: %w", err)
	}
	return nil
}