package main

import (
	"log/slog"
	"os"

	"GustavoCesarSantos/checkly-api/internal/infra/database"
	http "GustavoCesarSantos/checkly-api/internal/infra/http"

	"github.com/joho/godotenv"
)

func main() {
	// if os.Getenv("LOAD_ENV_FILE") == "true" {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}
	// }
	db, openDBErr := database.OpenDB()
	if openDBErr != nil {
		slog.Error(openDBErr.Error())
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("[API] Database connection pool established")
	err := http.Server(db)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
