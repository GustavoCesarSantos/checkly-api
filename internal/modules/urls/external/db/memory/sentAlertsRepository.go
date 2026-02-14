package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
)

type sentAlerts struct {
	IdempotencyKey string
	Status         domain.AlertStatus
	SentAt         *time.Time
}

type sentAlertsRepository struct {
	mu   sync.Mutex
	data map[string]*sentAlerts
	next int64
}

func NewSentAlertsRepository() db.ISentAlertsRepository {
	return &sentAlertsRepository{
		data: make(map[string]*sentAlerts),
		next: 1,
	}
}

func (r *sentAlertsRepository) Save(ctx context.Context, idempotencyKey string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	alert := sentAlerts{
		IdempotencyKey: idempotencyKey,
		Status:         domain.StatusPending,
	}
	r.data[alert.IdempotencyKey] = &alert
	return nil
}

func (r *sentAlertsRepository) Update(ctx context.Context, idempotencyKey string, status domain.AlertStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	alert, ok := r.data[idempotencyKey]
	if !ok {
		return errors.New("alert not found")
	}
	alert.Status = status
	sentAt := time.Now()
	alert.SentAt = &sentAt
	return nil
}
