package application

import (
	"context"
	"net/http"
	"time"
)

type CheckUrl struct {}

func NewCheckUrl() *CheckUrl {
	return &CheckUrl{}
}

type CheckUrlResult struct {
	IsSuccess bool
	StatusCode int
}

func (c *CheckUrl) Execute(url string) (CheckUrlResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CheckUrlResult{}, err
	}
	defer resp.Body.Close()
	if(resp.StatusCode >= 400) {
		return CheckUrlResult{
			IsSuccess: false,
			StatusCode: resp.StatusCode,
		}, nil
	}
	return CheckUrlResult{
		IsSuccess: true,
		StatusCode: resp.StatusCode,
	}, nil
}