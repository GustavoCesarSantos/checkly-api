package application

import (
	"errors"
	"fmt"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type ScheduleNextCheck struct{}

func NewScheduleNextCheck() *ScheduleNextCheck {
	return &ScheduleNextCheck{}
}

func (s *ScheduleNextCheck) Execute(url *domain.Url, now time.Time) error {
	if url.Status == domain.StatusHealthy {
		nextCheck := now.Add(time.Duration(url.Interval) * time.Minute)
		url.NextCheck = &nextCheck
		return nil
	}
	if url.Status == domain.StatusDegraded || url.Status == domain.StatusRecovering {
		nextCheck := now.Add(1 * time.Minute)
		url.NextCheck = &nextCheck
		return nil
	}
	if url.Status == domain.StatusDown {
		nextCheck := now.Add(url.Backoff())
		url.NextCheck = &nextCheck
		return nil
	}
	return fmt.Errorf("scheduleNextCheck: %w", errors.New("failed to schedule next check: unknown status"))
}
