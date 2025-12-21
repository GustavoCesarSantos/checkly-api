package main

import (
	"log/slog"
	"os"

	"GustavoCesarSantos/checkly-api/internal/infra/database"
	http "GustavoCesarSantos/checkly-api/internal/infra/http"

	"github.com/joho/godotenv"
)

// @title Checkly API
// @version 1.0
// @description API responsável por monitorar URLs e avaliar sua disponibilidade.
// @termsOfService http://swagger.io/terms/

// @contact.name Gustavo Cesar Santos
// @contact.email gustavocs789@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1
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
