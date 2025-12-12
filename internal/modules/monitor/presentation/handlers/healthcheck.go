package monitor

import (
	"net/http"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/monitor/presentation/dtos"
	"GustavoCesarSantos/checkly-api/internal/shared/configs"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type Healthcheck struct {}

func NewHealthcheck() *Healthcheck {
    return &Healthcheck{}
}

type HealthCheckEnvelope struct {
	HealthCheck dtos.HealthCheckResponse `json:"health_check"`
}

func (hc *Healthcheck) Handle(w http.ResponseWriter, r *http.Request) {
	metadataErr := utils.Envelope{
		"file": "healthcheck.go",
		"func": "healthcheck.Handle",
		"line": 0,
	}
	serverConfig := configs.LoadServerConfig()
	response := dtos.NewHealthCheckResponse("available", serverConfig.Env, time.Now().UTC().Format(time.RFC3339))
	err := utils.WriteJSON(w, http.StatusOK, utils.Envelope{"health_check": response}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err, metadataErr)
	}
}