package dtos

import (
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type UpdateUrlRequest struct {
	NextCheck		*time.Time			`json:"next_check"`
	RetryCount		*int				`json:"retry_count"`
	DownCount		*int				`json:"down_count"`
	StabilityCount	*int				`json:"stability_count"`
	Status			*domain.UrlStatus	`json:"status"`
}
