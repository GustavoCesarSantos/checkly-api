package dtos

type SystemInfo struct {
	Environment string `json:"environment" example:"development"`
	Time        string `json:"time" example:"2024-06-01T12:00:00Z"`
}

type HealthCheckResponse struct {
	Status     string     `json:"status" example:"available"`
	SystemInfo SystemInfo `json:"system_info"`
}

func NewHealthCheckResponse(status string, environment string, time string) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status: status,
		SystemInfo: SystemInfo{
			Environment: environment,
			Time:        time,
		},
	}
}
