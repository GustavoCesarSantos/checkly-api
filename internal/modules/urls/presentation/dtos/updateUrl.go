package dtos

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"time"
)

type UpdateUrlRequest struct {
	NextCheck *time.Time `json:"next_check"`
	RetryCount *int `json:"retry_count"`
	StabilityCount *int `json:"stability_count"`
	Status *domain.UrlStatus `json:"status"`
}