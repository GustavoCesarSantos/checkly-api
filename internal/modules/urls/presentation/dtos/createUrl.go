package dtos

import (
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
)

type CreateUrlRequest struct {
	Address    string            `json:"address" example:"https://api.meuservico.com/health"`
	Interval   int               `json:"interval_minutes" example:"5"`
	RetryLimit int               `json:"retry_limit" example:"3"`
	Contact    string            `json:"contact_email" example:"meu_contato@email.com"`
	Status     *domain.UrlStatus `json:"status"`
	NextCheck  *time.Time        `json:"next_check"`
}

type CreateUrlResponse struct {
	ID string `json:"id" example:"1"`
}

func NewCreateUrlResponse(externalId string) *CreateUrlResponse {
	return &CreateUrlResponse{
		ID: externalId,
	}
}
