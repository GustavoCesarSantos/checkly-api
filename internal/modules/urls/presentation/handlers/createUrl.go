package urls

import (
	"net/http"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type CreateUrl struct {
	urlRepository db.IUrlRepository
}

func NewCreateUrl(urlRepository db.IUrlRepository) *CreateUrl {
	return &CreateUrl{
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
	url := domain.NewUrl(0, int64(input.Interval), int64(input.RetryLimit), input.Contact)
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