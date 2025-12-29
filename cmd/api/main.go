package main

import (
	"os"
	"sync"

	"GustavoCesarSantos/checkly-api/internal/infra/database"
	http "GustavoCesarSantos/checkly-api/internal/infra/http"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"

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

// @BasePath /v1
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
	var wg sync.WaitGroup
	err := http.Server(db, &wg)
	if err != nil {
		logger.Error(
			"failed to start HTTP server",
			"main.go",
			"main",
			err,
		)
		os.Exit(1)
	}
}
