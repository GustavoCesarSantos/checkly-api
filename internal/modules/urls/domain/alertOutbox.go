package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type AlertStatus int

const (
	StatusPending AlertStatus = 10
	StatusSent    AlertStatus = 20
	StatusFailed  AlertStatus = 30
)

type Payload struct {
	Url   string `json:"url"`
	Email string `json:"email"`
}

func (p Payload) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Payload) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type %T for Alert Payload", value)
	}
	return json.Unmarshal(bytes, p)
}

type AlertOutbox struct {
	ID           int64
	UrlId        int64
	Payload      Payload
	SentAt       *time.Time
	ProcessingAt *time.Time
	RetryCount   int
	NextRetryAt  time.Time
	LockedAt     *time.Time
	LockedBy     string
	CreatedAt    time.Time
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
