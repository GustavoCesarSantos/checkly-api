package utils_urls

type UrlStatus int

const (
    StatusHealthy  UrlStatus = 10
    StatusDegraded UrlStatus = 20
    StatusRecovering UrlStatus = 25
    StatusDown     UrlStatus = 30
    StatusNotified UrlStatus = 40
)