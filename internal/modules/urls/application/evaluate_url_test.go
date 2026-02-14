package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"testing"
)

func TestEvaluateUrl_HealthyToDegraded_WhenHttpFails(t *testing.T) {
	url := &domain.Url{
		Status: domain.StatusHealthy,
	}
	sut := NewEvaluateUrl()
	sut.Execute(url, false)
	if url.Status != domain.StatusDegraded {
		t.Errorf("expected status Degraded, got %v", url.Status)
	}
}

func TestEvaluateUrl_DegradedToRecovering_WhenHttpOk(t *testing.T) {
	url := &domain.Url{
		Status:     domain.StatusDegraded,
		RetryCount: 5,
	}
	sut := NewEvaluateUrl()
	sut.Execute(url, true)
	if url.Status != domain.StatusRecovering {
		t.Errorf("expected status Recovering, got %v", url.Status)
	}
	if url.RetryCount != 0 {
		t.Errorf("expected RetryCount reset to 0, got %d", url.RetryCount)
	}
}

func TestEvaluateUrl_RecoveringToHealthy_AfterThreeSuccesses(t *testing.T) {
	url := &domain.Url{
		Status:         domain.StatusRecovering,
		StabilityCount: 2,
	}
	sut := NewEvaluateUrl()
	sut.Execute(url, true)
	if url.Status != domain.StatusHealthy {
		t.Errorf("expected status Healthy, got %v", url.Status)
	}
	if url.StabilityCount != 0 {
		t.Errorf("expected StabilityCount reset to 0, got %d", url.StabilityCount)
	}
}

func TestEvaluateUrl_RecoveringToDegraded_WhenHttpFails(t *testing.T) {
	url := &domain.Url{
		Status:         domain.StatusRecovering,
		StabilityCount: 2,
	}
	sut := NewEvaluateUrl()
	sut.Execute(url, false)
	if url.Status != domain.StatusDegraded {
		t.Errorf("expected status Degraded, got %v", url.Status)
	}
	if url.StabilityCount != 0 {
		t.Errorf("expected StabilityCount reset to 0, got %d", url.StabilityCount)
	}
}

func TestEvaluateUrl_DegradedToDown_WhenRetryLimitReached(t *testing.T) {
	url := &domain.Url{
		Status:     domain.StatusDegraded,
		RetryCount: 3,
		RetryLimit: 3,
	}
	sut := NewEvaluateUrl()
	sut.Execute(url, false)
	if url.Status != domain.StatusDown {
		t.Errorf("expected status Down, got %v", url.Status)
	}
	if url.RetryCount != 0 {
		t.Errorf("expected RetryCount reset to 0, got %d", url.RetryCount)
	}
}

func TestEvaluateUrl_DegradedRetryIncrement_WhenBelowLimit(t *testing.T) {
	url := &domain.Url{
		Status:     domain.StatusDegraded,
		RetryCount: 1,
		RetryLimit: 3,
	}
	sut := NewEvaluateUrl()
	sut.Execute(url, false)
	if url.RetryCount != 2 {
		t.Errorf("expected RetryCount=2, got %d", url.RetryCount)
	}
}
