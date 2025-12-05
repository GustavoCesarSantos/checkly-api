package domain

import "time"

type Url struct {
    ID int64
    Url string
	Interval int
    RetryLimit int
	Contact string
    CreatedAt time.Time
    UpdatedAt *time.Time
}

func NewUrl(
    id int64,
    url string,
	interval int,
    retryLimit int,
	contact string,
) *Url {
    return &Url{
        ID: id,
        Url: url,
		Interval: interval,
        RetryLimit: retryLimit,
        Contact: contact,
        CreatedAt: time.Now(),
    }
}