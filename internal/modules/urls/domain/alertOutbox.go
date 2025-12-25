package domain

import "time"

type Payload struct {
	Email   string `json:"email"`
}

type AlertOutbox struct {
	ID             int64
	UrlId          int64
	Payload        Payload
	IdempotencyKey string
	SentAt         *time.Time
	ProcessingAt   *time.Time
	RetryCount     int
	NextRetryAt    time.Time
	LockedAt	   *time.Time
	LockedBy	   string
	CreatedAt      time.Time
}

func NewAlertOutbox(
	urlId int64,
	payload Payload,
) *AlertOutbox {
	return &AlertOutbox{
		UrlId:   urlId,
		Payload: payload,
	}
}