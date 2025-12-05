package urls

import (
	"net/http"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type CreateUrl struct {
	checkUrl application.CheckUrl
	urlRepository db.IUrlRepository
}

func NewCreateUrl(checkUrl application.CheckUrl, urlRepository db.IUrlRepository) *CreateUrl {
	return &CreateUrl{
		checkUrl,
		urlRepository,
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
	if(!checkResult.IsSuccess) {
		utils.ServerErrorResponse(w, r, utils.ErrFailedCheckUrl, metadataErr)
		return
	}
	url := domain.NewUrl(
		0, 
		input.Url, 
		input.Interval, 
		input.RetryLimit, 
		input.Contact,
	)
	saveErr := cu.urlRepository.Save(url)
	if saveErr != nil {
		utils.ServerErrorResponse(w, r, readErr, metadataErr)
		return
	}
	response := dtos.NewCreateUrlResponse(url.ID)
	err := utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"url": response}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err, metadataErr)
	}
}