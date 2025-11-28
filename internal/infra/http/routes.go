package http

import (
	"database/sql"
	"net/http"

	"GustavoCesarSantos/checkly-api/internal/infra/http/middleware"
	monitor "GustavoCesarSantos/checkly-api/internal/modules/monitor/presentation/handlers"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
)


func routes(mux *http.ServeMux, sqlDB *sql.DB) http.Handler {
	urlRepository := db.NewUrlRepository(sqlDB)
	
	createUrl := urls.NewCreateUrl(urlRepository)
	healthcheck := monitor.NewHealthcheck()

	mux.Handle("GET /v1/health", http.HandlerFunc(healthcheck.Handle))
	mux.Handle("POST /v1/urls", http.HandlerFunc(createUrl.Handle))
	return middleware.RecoverPanic(middleware.EnableCORS(mux))
}