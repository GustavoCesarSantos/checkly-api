package domain

import (
	"time"
)

type UrlStatus int

const (
	StatusHealthy    UrlStatus = 10
	StatusDegraded   UrlStatus = 20
	StatusRecovering UrlStatus = 25
	StatusDown       UrlStatus = 30
)

type Url struct {
	ID             int64
	ExternalID     string
	Address        string
	Interval       int
	RetryLimit     int
	RetryCount     int
	DownCount      int
	WentDownNow    bool
	StabilityCount int
	Contact        string
	NextCheck      *time.Time
	Status         UrlStatus
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

func NewUrl(
	address string,
	interval int,
	retryLimit int,
	contact string,
	status UrlStatus,
	nextCheck time.Time,
) *Url {
	return &Url{
		Address:    address,
		Interval:   interval,
		RetryLimit: retryLimit,
		Contact:    contact,
		Status:     status,
		NextCheck:  &nextCheck,
	}
}

func (u Url) Backoff() time.Duration {
	steps := []time.Duration{
		1 * time.Minute,
		2 * time.Minute,
		5 * time.Minute,
		10 * time.Minute,
		30 * time.Minute,
	}
	if u.DownCount <= 0 {
		return steps[0]
	}
	if u.DownCount > len(steps) {
		return steps[len(steps)-1]
	}
	return steps[u.DownCount-1]
}
