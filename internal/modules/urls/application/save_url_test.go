package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/memory"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"context"
	"errors"
	"testing"
	"time"
)

type failingRepository struct{}

var errSave = errors.New("save failed")

func (f *failingRepository) Save(url *domain.Url) error {
	return errSave
}

func (f *failingRepository) FindAllByNextCheck(_ context.Context, _ time.Time) ([]domain.Url, error) {
	panic("not used")
}

func (f *failingRepository) Update(_ context.Context, _ int64, _ db.UpdateUrlParams) error {
	panic("not used")
}

func TestSaveUrl_Execute_Healthy(t *testing.T) {
	repo := memory.NewUrlRepository()
	sut := NewSaveUrl(repo)
	input := dtos.CreateUrlRequest{
		Address:    "https://example.com",
		Interval:   5,
		RetryLimit: 3,
		Contact:    "admin@example.com",
	}
	now := time.Now()
	url, err := sut.Execute(input, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url == nil {
		t.Fatal("expected url, got nil")
	}
	if url.Status != domain.StatusHealthy {
		t.Errorf("expected status Healthy, got %v", url.Status)
	}
	if url.ID == 0 {
		t.Errorf("expected url to have an ID")
	}
	expectedMin := now.Add(5 * time.Minute)
	expectedMax := expectedMin.Add(2 * time.Second)
	if url.NextCheck.Before(expectedMin) || url.NextCheck.After(expectedMax) {
		t.Errorf("NextCheck out of expected range, got %v", url.NextCheck)
	}
}

func TestSaveUrl_Execute_Degraded(t *testing.T) {
	repo := memory.NewUrlRepository()
	sut := NewSaveUrl(repo)
	input := dtos.CreateUrlRequest{
		Address:    "https://example.com",
		Interval:   10,
		RetryLimit: 5,
		Contact:    "admin@example.com",
	}
	now := time.Now()
	url, err := sut.Execute(input, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url.Status != domain.StatusDegraded {
		t.Errorf("expected status Degraded, got %v", url.Status)
	}
	expectedMin := now.Add(1 * time.Minute)
	expectedMax := expectedMin.Add(2 * time.Second)
	if url.NextCheck.Before(expectedMin) || url.NextCheck.After(expectedMax) {
		t.Errorf("NextCheck out of expected range, got %v", url.NextCheck)
	}
}

func TestSaveUrl_Execute_WhenRepositoryFails(t *testing.T) {
	repo := &failingRepository{}
	sut := NewSaveUrl(repo)
	input := dtos.CreateUrlRequest{
		Address:    "https://example.com",
		Interval:   5,
		RetryLimit: 3,
		Contact:    "admin@example.com",
	}
	url, err := sut.Execute(input, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if url != nil {
		t.Errorf("expected nil url, got %v", url)
	}
}
