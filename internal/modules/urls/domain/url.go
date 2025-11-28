package domain

import "time"

type Url struct {
    ID int64
	Interval int64
    RetryLimit int64
	Contact string
    CreatedAt time.Time
    UpdatedAt *time.Time
}

func NewUrl(
    id int64,
	interval int64,
    retryLimit int64,
	contact string,
) *Url {
    return &Url{
        ID: id,
		Interval: interval,
        RetryLimit: retryLimit,
        Contact: contact,
        CreatedAt: time.Now(),
    }
}