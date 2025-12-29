package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"GustavoCesarSantos/checkly-api/internal/infra/database"
	"GustavoCesarSantos/checkly-api/internal/infra/worker"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("LOAD_ENV_FILE") == "true" {
		if err := godotenv.Load(); err != nil {
			logger.Error(
				"failed to load .env file",
				"main.go",
				"main",
				err,
			)
			os.Exit(1)
		}
	}
	db, openDBErr := database.OpenDB()
	if openDBErr != nil {
		logger.Error(
			"failed to open database",
			"main.go",
			"main",
			openDBErr,
		)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info(
		"database connection pool established",
		"main.go",
		"main",
	)
	mw := worker.NewMonitorWorker(db, 5)
	nw := worker.NewNotifyWorker(db, 5)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	mw.Start(ctx)
	nw.Start(ctx)
}
