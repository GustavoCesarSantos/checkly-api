package configs

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Port int
	Env  string
	Cors struct {
		TrustedOrigins []string
	}
	Jwt struct {
		Secret string
	}
}

func LoadServerConfig() ServerConfig {
	port, portErr := strconv.Atoi(GetEnv("PORT", "9000"))
	if portErr != nil {
		slog.Error(portErr.Error())
		os.Exit(1)
	}
	return ServerConfig{
		Port: port,
		Env:  GetEnv("ENV", "development"),
		Cors: struct {
			TrustedOrigins []string
		}{
			TrustedOrigins: strings.Fields(GetEnv("TRUSTED_ORIGINS", "*")),
		},
		Jwt: struct {
			Secret string
		}{
			Secret: GetEnv("JWT_SECRET", "defaultsecret"),
		},
	}
}
