package application

import (
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type ScheduleNextCheck struct{}

func NewScheduleNextCheck() *ScheduleNextCheck {
	return &ScheduleNextCheck{}
}

func (s *ScheduleNextCheck) Execute(url *domain.Url, httpOK bool, now time.Time) {
	nextCheck := now.Add(time.Duration(url.Interval) * time.Minute)
	if !httpOK {
		nextCheck = now.Add(1 * time.Minute)
	}
	url.NextCheck = &nextCheck
}