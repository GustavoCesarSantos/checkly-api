package application

import (
	"fmt"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"GustavoCesarSantos/checkly-api/internal/shared/mailer"
)

type SendEmail struct {
	mailer mailer.Mailer
}

func NewSendEmail(mailer mailer.Mailer) *SendEmail {
	return &SendEmail{mailer: mailer}
}

func (s *SendEmail) Execute(payload domain.Payload) error {
	data := map[string]any{
		"url":          payload.Url,
	}
	err := s.mailer.Send(payload.Email, "alert_url_down.tmpl", data)
	if err != nil {
		return fmt.Errorf("sendEmail: %w", err)
	}
	return nil
}