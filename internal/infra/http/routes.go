package http

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "GustavoCesarSantos/checkly-api/docs"
	"GustavoCesarSantos/checkly-api/internal/infra/http/middleware"
	monitor "GustavoCesarSantos/checkly-api/internal/modules/monitor/presentation/handlers"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

func routes(sqlDB *sql.DB) http.Handler {
	urlRepository := db.NewUrlRepository(sqlDB)
	
	checkUrl := application.NewCheckUrl()
	saveUrl := application.NewSaveUrl(urlRepository)
	
	createUrl := urls.NewCreateUrl(checkUrl, saveUrl)
	healthcheck := monitor.NewHealthcheck()
	
	router := httprouter.New()

	metadataErr := utils.Envelope{
		"file": "routes.go",
		"func": "routes.routes",
		"line": 0,
	}

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadataErr["line"] = 42
		utils.NotFoundResponse(w, r, metadataErr)
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadataErr["line"] = 47
		utils.MethodNotAllowedResponse(w, r, metadataErr)
	})

	router.Handler(http.MethodGet, "/v1/docs/*filepath", httpSwagger.WrapHandler)
	router.HandlerFunc(http.MethodGet, "/v1/health", healthcheck.Handle)
	router.HandlerFunc(http.MethodPost, "/v1/urls", createUrl.Handle)
	return middleware.RecoverPanic(middleware.EnableCORS(router))
}
