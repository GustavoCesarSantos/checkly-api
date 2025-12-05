package urls

import (
	"net/http"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	utils_urls "GustavoCesarSantos/checkly-api/internal/modules/urls/utils"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type CreateUrl struct {
	checkUrl *application.CheckUrl
	saveUrl *application.SaveUrl
}

func NewCreateUrl(checkUrl *application.CheckUrl, saveUrl *application.SaveUrl) *CreateUrl {
	return &CreateUrl{
		checkUrl,
		saveUrl,
	}
}

type CreateUrlEnvelop struct {
	CreateUrl dtos.CreateUrlResponse `json:"url"`
}

func (cu *CreateUrl) Handle(w http.ResponseWriter, r *http.Request) {
	metadataErr := utils.Envelope{
		"file": "createUrl.go",
		"func": "createUrl.Handle",
		"line": 0,
	}
	var input dtos.CreateUrlRequest
	readErr := utils.ReadJSON(w, r, &input)
	if readErr != nil {
		utils.BadRequestResponse(w, r, readErr, metadataErr)
		return
	}
	checkResult, checkErr := cu.checkUrl.Execute(input.Url)
	if(checkErr != nil) {
		utils.ServerErrorResponse(w, r, utils.ErrFailedCheckUrl, metadataErr)
		return
	}
	status := utils_urls.StatusHealthy
	nextCheck := time.Now().Add(time.Duration(input.Interval) * time.Minute)
	if(!checkResult.IsSuccess) {
		status = utils_urls.StatusDegraded
		nextCheck = time.Now().Add(time.Minute)
	}
	input.Status = &status
	input.NextCheck = &nextCheck
	url, saveErr := cu.saveUrl.Execute(input)
	if saveErr != nil {
		utils.ServerErrorResponse(w, r, saveErr, metadataErr)
		return
	}
	response := dtos.NewCreateUrlResponse(url.ID)
	err := utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"url": response}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err, metadataErr)
	}
}