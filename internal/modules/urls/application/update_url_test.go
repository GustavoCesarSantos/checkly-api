package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/memory"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

type failingUpdateRepository struct{}

var errUpdate = errors.New("update failed")

func (f *failingUpdateRepository) FindAllByNextCheck(context.Context, time.Time) ([]domain.Url, error) {
	panic("not used")
}

func (f *failingUpdateRepository) Save(ctx context.Context, url *domain.Url) error {
	panic("not used")
}

func (f *failingUpdateRepository) Update(
	ctx context.Context,
	_ int64,
	_ db.UpdateUrlParams,
) error {
	return errUpdate
}

func (f *failingUpdateRepository) UpdateToNotified(_ context.Context, _ int64) error {
	panic("not used")
}

func TestUpdateUrl_Execute_UpdateFields(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := memory.NewUrlRepository()
	sut := NewUpdateUrl(repo)

	url := domain.NewUrl(
		"https://example.com",
		5,
		3,
		"admin@example.com",
		domain.StatusHealthy,
		time.Now(),
	)

	if err := repo.Save(ctx, url); err != nil {
		t.Fatalf("failed to save url: %v", err)
	}

	nextCheck := time.Now().Add(10 * time.Minute)
	retryCount := 2
	stability := 1
	status := domain.StatusDegraded

	input := dtos.UpdateUrlRequest{
		NextCheck:      &nextCheck,
		RetryCount:     &retryCount,
		StabilityCount: &stability,
		Status:         &status,
	}

	// Act
	err := sut.Execute(ctx, url.ID, input)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := repo.FindAllByNextCheck(ctx, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("failed to query urls: %v", err)
	}

	var found *domain.Url
	for i := range updated {
		if updated[i].ID == url.ID {
			found = &updated[i]
			break
		}
	}

	if found == nil {
		t.Fatal("updated url not found")
	}

	if found.Status != domain.StatusDegraded {
		t.Errorf("expected status Degraded, got %v", found.Status)
	}

	if found.RetryCount != retryCount {
		t.Errorf("expected RetryCount=%d, got %d", retryCount, found.RetryCount)
	}

	if found.StabilityCount != stability {
		t.Errorf("expected StabilityCount=%d, got %d", stability, found.StabilityCount)
	}

	if !found.NextCheck.Equal(nextCheck) {
		t.Errorf("expected NextCheck=%v, got %v", nextCheck, found.NextCheck)
	}
}

func TestUpdateUrl_Execute_DoesNotOverrideNilFields(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := memory.NewUrlRepository()
	sut := NewUpdateUrl(repo)

	originalNextCheck := time.Now().Add(5 * time.Minute)

	url := domain.NewUrl(
		"https://example.com",
		5,
		3,
		"admin@example.com",
		domain.StatusHealthy,
		originalNextCheck,
	)

	if err := repo.Save(ctx, url); err != nil {
		t.Fatalf("failed to save url: %v", err)
	}

	newStatus := domain.StatusDegraded

	input := dtos.UpdateUrlRequest{
		Status: &newStatus,
		// outros campos nil
	}

	// Act
	err := sut.Execute(ctx, url.ID, input)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := repo.FindAllByNextCheck(ctx, time.Now().Add(1*time.Hour))
	found := updated[0]

	if found.Status != domain.StatusDegraded {
		t.Errorf("expected status Degraded, got %v", found.Status)
	}

	if !found.NextCheck.Equal(originalNextCheck) {
		t.Errorf("expected NextCheck unchanged, got %v", found.NextCheck)
	}
}

func TestUpdateUrl_Execute_WhenRepositoryFails(t *testing.T) {
	sut := NewUpdateUrl(&failingUpdateRepository{})
	input := dtos.UpdateUrlRequest{}
	err := sut.Execute(context.Background(), 1, input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
