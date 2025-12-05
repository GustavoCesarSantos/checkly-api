package application

import "net/http"

type CheckUrl struct {}

func NewCheckUrl() *CheckUrl {
	return &CheckUrl{}
}

type CheckUrlResult struct {
	IsSuccess bool
	StatusCode int
}

func (cu *CheckUrl) Execute(url string) (CheckUrlResult, error) {
	resp, err := http.Get(url)
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