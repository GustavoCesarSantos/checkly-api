package dtos

type CreateUrlRequest struct {
	Url string `json:"url" example:"ahttps://api.meuservico.com/health"`
	Interval int `json:"interval_minutes" example:"5"`
	RetryLimit int `json:"retry_limit" example:"3"`
	Contact string `json:"contact_email" example:"meu_contato@email.com"`
}

type CreateUrlResponse struct {
	ID int64 `json:"id" example:"1"`
}

func NewCreateUrlResponse(id int64) *CreateUrlResponse {
	return &CreateUrlResponse{
		ID: id,
	}
}