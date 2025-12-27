package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"GustavoCesarSantos/checkly-api/internal/shared/mailer"
	"log/slog"
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
		slog.Warn("failed to send email", "to", payload.Email, "error", err.Error())
		return err
	}
	slog.Info("email sent successfully", "to", payload.Email)
	return nil
}