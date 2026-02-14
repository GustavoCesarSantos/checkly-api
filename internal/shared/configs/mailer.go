package configs

import (
	"log/slog"
	"os"
	"strconv"
)

type MailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func LoadMailerConfig() MailerConfig {
	port, portErr := strconv.Atoi(GetEnv("MAILER_PORT", "25"))
	if portErr != nil {
		slog.Error(portErr.Error())
		os.Exit(1)
	}
	return MailerConfig{
		Host:     GetEnv("MAILER_HOST", "localhost"),
		Port:     port,
		Username: GetEnv("MAILER_USERNAME", ""),
		Password: GetEnv("MAILER_PASSWORD", ""),
		Sender:   GetEnv("MAILER_SENDER", ""),
	}
}
