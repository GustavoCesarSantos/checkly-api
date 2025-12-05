package domain

import (
	utils_urls "GustavoCesarSantos/checkly-api/internal/modules/urls/utils"
	"time"
)

type Url struct {
    ID int64
    Url string
	Interval int
    RetryLimit int
    RetryCount int
    StabilityCount int
	Contact string
    NextCheck *time.Time
    Status utils_urls.UrlStatus
    CreatedAt time.Time
    UpdatedAt *time.Time
}

func NewUrl(
    url string,
	interval int,
    retryLimit int,
	contact string,
    status utils_urls.UrlStatus,
    nextCheck time.Time,
) *Url {
    return &Url{
        Url: url,
		Interval: interval,
        RetryLimit: retryLimit,
        Contact: contact,
        Status: status,
        NextCheck: &nextCheck,
    }
}