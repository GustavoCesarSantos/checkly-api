package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"testing"
	"time"
)

func TestScheduleNextCheck_NextCheck_UsesInterval_WhenHttpOk(t *testing.T) {
	now := time.Now()
	url := &domain.Url{
		Status:   domain.StatusHealthy,
		Interval: 5,
	}
	sut := NewScheduleNextCheck()
	sut.Execute(url, true, now)
	if url.NextCheck == nil {
		t.Fatal("expected NextCheck to be set")
	}
	expectedMin := now.Add(5 * time.Minute)
	expectedMax := expectedMin.Add(2 * time.Second)
	if url.NextCheck.Before(expectedMin) || url.NextCheck.After(expectedMax) {
		t.Errorf("NextCheck not in expected range, got %v", url.NextCheck)
	}
}

func TestScheduleNextCheck_NextCheck_OneMinute_WhenHttpFails(t *testing.T) {
	now := time.Now()
	url := &domain.Url{
		Status:   domain.StatusHealthy,
		Interval: 10,
	}
	sut := NewScheduleNextCheck()
	sut.Execute(url, false, now)
	if url.NextCheck == nil {
		t.Fatal("expected NextCheck to be set")
	}
	expectedMin := now.Add(1 * time.Minute)
	expectedMax := expectedMin.Add(2 * time.Second)
	if url.NextCheck.Before(expectedMin) || url.NextCheck.After(expectedMax) {
		t.Errorf("NextCheck not in expected range, got %v", url.NextCheck)
	}
}
