package dtos

type SystemInfo struct {
	Environment string `json:"environment" example:"development"`
}

type HealthCheckResponse struct {
	Status string `json:"status" example:"available"`
	SystemInfo SystemInfo `json:"system_info"`
}

func NewHealthCheckResponse(status string, environment string) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status: status,
		SystemInfo: SystemInfo{
			Environment: environment,
		},
	}
}