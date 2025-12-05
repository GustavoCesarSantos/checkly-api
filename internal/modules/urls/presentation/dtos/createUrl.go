package dtos

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"time"
)

type CreateUrlRequest struct {
	Url string `json:"url" example:"ahttps://api.meuservico.com/health"`
	Interval int `json:"interval_minutes" example:"5"`
	RetryLimit int `json:"retry_limit" example:"3"`
	Contact string `json:"contact_email" example:"meu_contato@email.com"`
	Status *domain.UrlStatus  `json:"status"`
	NextCheck *time.Time `json:"next_check"`
}

type CreateUrlResponse struct {
	ID int64 `json:"id" example:"1"`
}

func NewCreateUrlResponse(id int64) *CreateUrlResponse {
	return &CreateUrlResponse{
		ID: id,
	}
}