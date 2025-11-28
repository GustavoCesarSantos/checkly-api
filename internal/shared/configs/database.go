package configs

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig struct {
	Dsn string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime time.Duration
}

func LoadDatabaseConfig() DatabaseConfig {
	maxOpenConns, maxOpenConnsErr := strconv.Atoi(GetEnv("DB_MAX_OPEN_CONNS", "25"))
	if(maxOpenConnsErr != nil) {
		slog.Error(maxOpenConnsErr.Error())
		os.Exit(1)
	}
	maxIdleConns, maxIdleConnsErr := strconv.Atoi(GetEnv("DB_MAX_IDLE_CONNS", "25"))
	if(maxIdleConnsErr != nil) {
		slog.Error(maxIdleConnsErr.Error())
		os.Exit(1)
	}
	maxIdleTime, maxIdleTimeErr := time.ParseDuration(GetEnv("DB_MAX_IDLE_TIME", "15m"))
	if(maxIdleTimeErr != nil) {
		slog.Error(maxIdleTimeErr.Error())
		os.Exit(1)
	}
	return DatabaseConfig {
		Dsn: GetEnv("DB_DSN", "postgres://usuario:senha@localhost:5432/testDb?sslmode=disable"),
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
		MaxIdleTime: maxIdleTime,
	}
}