package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
)

type urlRepository struct {
	mu   sync.Mutex
	data map[int64]*domain.Url
	next int64
}

func NewUrlRepository() db.IUrlRepository {
	return &urlRepository{
		data: make(map[int64]*domain.Url),
		next: 1,
	}
}

func (r *urlRepository) FindAllByNextCheck(
	ctx context.Context,
	nextCheck time.Time,
) ([]domain.Url, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []domain.Url
	for _, url := range r.data {
		if url.NextCheck == nil {
			continue
		}
		if url.NextCheck.Before(nextCheck) || url.NextCheck.Equal(nextCheck) {
			result = append(result, *url)
		}
	}
	return result, nil
}

func (r *urlRepository) Save(url *domain.Url) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	url.ID = r.next
	r.next++
	r.data[url.ID] = url
	return nil
}

func (r *urlRepository) Update(
	ctx context.Context,
	urlId int64,
	params db.UpdateUrlParams,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	url, ok := r.data[urlId]
	if !ok {
		return errors.New("url not found")
	}
	if params.NextCheck != nil {
		url.NextCheck = params.NextCheck
	}
	if params.RetryCount != nil {
		url.RetryCount = *params.RetryCount
	}
	if params.DownCount != nil {
		url.DownCount = *params.DownCount
	}
	if params.StabilityCount != nil {
		url.StabilityCount = *params.StabilityCount
	}
	if params.Status != nil {
		url.Status = *params.Status
	}
	return nil
}
