package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"GustavoCesarSantos/checkly-api/internal/infra/database"
	"GustavoCesarSantos/checkly-api/internal/infra/worker"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("LOAD_ENV_FILE") == "true" {
		if err := godotenv.Load(); err != nil {
			slog.Error("failed to load .env file", "error", err)
			os.Exit(1)
		}
	}
	db, openDBErr := database.OpenDB()
	if openDBErr != nil {
		slog.Error(openDBErr.Error())
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database connection pool established")
	mw := worker.NewMonitorWorker(db, 5)
	nw := worker.NewNotifyWorker(db, 5)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	mw.Start(ctx)
	nw.Start(ctx)
}
