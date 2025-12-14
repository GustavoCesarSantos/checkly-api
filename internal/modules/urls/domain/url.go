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
	StatusNotified   UrlStatus = 40
)

type Url struct {
	ID             int64
	ExternalID     string
	Address        string
	Interval       int
	RetryLimit     int
	RetryCount     int
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
