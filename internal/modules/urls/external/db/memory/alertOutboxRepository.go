package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
)

type alertOutboxRepository struct {
	mu   sync.Mutex
	data map[int64]*domain.AlertOutbox
	next int64
}

func NewAlertOutboxRepository() db.IAlertOutboxRepository {
	return &alertOutboxRepository{
		data: make(map[int64]*domain.AlertOutbox),
		next: 1,
	}
}

func (r *alertOutboxRepository) FindAllPendingAlerts(ctx context.Context, limit int) ([]domain.AlertOutbox, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []domain.AlertOutbox
	for _, alertOutbox := range r.data {
		if alertOutbox.NextRetryAt.Before(time.Now()) {
			result = append(result, *alertOutbox)
		}
	}
	return result, nil
}

func (r *alertOutboxRepository) Save(ctx context.Context, alertOutbox *domain.AlertOutbox) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	alertOutbox.ID = r.next
	r.next++
	r.data[alertOutbox.ID] = alertOutbox
	return nil
}

func (r *alertOutboxRepository) Update(ctx context.Context, alertId int64, sentAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	alertOutbox, ok := r.data[alertId]
	if !ok {
		return errors.New("alert not found")
	}
	alertOutbox.SentAt = &sentAt
	return nil
}

func (r *alertOutboxRepository) UpdateRetryInfo(ctx context.Context, alertId int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	alertOutbox, ok := r.data[alertId]
	if !ok {
		return errors.New("alert not found")
	}
	alertOutbox.RetryCount++
	alertOutbox.NextRetryAt = time.Now().Add(time.Duration(alertOutbox.RetryCount) * time.Minute)
	return nil
}
