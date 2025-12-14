package urls

import (
	"net/http"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
	"GustavoCesarSantos/checkly-api/internal/shared/validator"
)

type CreateUrl struct {
	checkUrl *application.CheckUrl
	saveUrl  *application.SaveUrl
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

func (c *CreateUrl) Handle(w http.ResponseWriter, r *http.Request) {
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
	v := c.ValidateInput(input)
	if !v.Valid() {
		utils.FailedValidationResponse(w, r, v.Errors, metadataErr)
		return
	}
	checkResult, checkErr := c.checkUrl.Execute(input.Address)
	if checkErr != nil {
		utils.ServerErrorResponse(w, r, utils.ErrFailedCheckUrl, metadataErr)
		return
	}
	url, saveErr := c.saveUrl.Execute(input, checkResult.IsSuccess)
	if saveErr != nil {
		utils.ServerErrorResponse(w, r, saveErr, metadataErr)
		return
	}
	response := dtos.NewCreateUrlResponse(url.ExternalID)
	err := utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"url": response}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err, metadataErr)
	}
}

func (c *CreateUrl) ValidateInput(input dtos.CreateUrlRequest) *validator.Validator {
	v := validator.NewValidator()
	v.Check(input.Address != "", "address", "must be provided")
	v.Check(validator.Matches(input.Address, validator.UrlRX), "address", "must be a valid URL")
	v.Check(input.Interval != 0, "interval_minutes", "must be provided")
	v.Check(input.Interval < 1, "interval_minutes", "must be greater than 0")
	v.Check(input.Interval > 60, "interval_minutes", "must be less than 60")
	v.Check(input.RetryLimit != 0, "retry_limit", "must be provided")
	v.Check(input.RetryLimit < 1, "retry_limit", "must be greater than 0")
	v.Check(input.RetryLimit > 10, "retry_limit", "must be less than 10")
	v.Check(input.Contact != "", "contact_email", "must be provided")
	v.Check(validator.Matches(input.Contact, validator.EmailRX), "contact_email", "must be a valid email address")
	return v
}
