package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"time"
)

type EvaluateUrl struct {}

func NewEvaluateUrl() *EvaluateUrl {
	return &EvaluateUrl{}
}

func (e *EvaluateUrl) Execute(url *domain.Url, httpOK bool) {
	switch {
		case url.Status == domain.StatusHealthy && !httpOK:
			url.Status = domain.StatusDegraded
			return
		case url.Status == domain.StatusDegraded && httpOK:
			url.Status = domain.StatusRecovering
			url.RetryCount = 0
			return
		case url.Status == domain.StatusRecovering && httpOK:
			url.StabilityCount++ 
			if(url.StabilityCount >= 3) {
				url.Status = domain.StatusHealthy
				url.StabilityCount = 0 
				return
			}
		case url.Status == domain.StatusRecovering && !httpOK:
			url.Status = domain.StatusDegraded
			url.StabilityCount = 0
			return
		case url.Status == domain.StatusDegraded && !httpOK:
			if(url.RetryCount >= url.RetryLimit) {
				url.Status = domain.StatusDown
				url.RetryCount = 0
				return
			}
			url.RetryCount++
			return
	}
	nextCheck := time.Now().Add(time.Duration(url.Interval) * time.Minute)
	if(!httpOK) {
		nextCheck = time.Now().Add(1 * time.Minute)
	}
	url.NextCheck = &nextCheck
}