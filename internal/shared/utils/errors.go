package utils

import (
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMissingJSONValue                = errors.New("BODY MUST CONTAIN A JSON VALUE")
	ErrMissingOrInvalidLimitQueryParam = errors.New("MISSING OR INVALID LIMIT QUERY PARAM")
	ErrInvalidLimitQueryParam          = errors.New("INVALID LAST ID QUERY PARAM")
	ErrEditConflict                    = errors.New("EDIT CONFLICT")
	ErrRecordNotFound                  = errors.New("RECORD NOT FOUND")
	ErrFailedCheckUrl                  = errors.New("FAILED TO CHECK URL")
)

type ErrorEnvelope struct {
	Error string `json:"error" example:"error message"`
}

type MetadataErr struct {
	Message string `json:"message" example:"detailed error message"`
	Who string `json:"who" example:"module or component that raised the error"`
	Where string `json:"where" example:"file and function where the error occurred"`
}

func logError(r *http.Request, err error, metadataErr MetadataErr) {
	var (
		message = metadataErr.Message
		who	= metadataErr.Who
		where = metadataErr.Where
		method = r.Method
		url    = r.URL.RequestURI()
	)
	logger.Error(
		message,
		who,
		where,
		err,
		"method", method, 
		"url", url,
	)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, err error, metadataErr MetadataErr) {
	logError(r, err, metadataErr)
	data := Envelope{"error": message}
	writeErr := WriteJSON(w, status, data, nil)
	if writeErr != nil {
		logger.Error(
			"Failed to write error response",
			"utils/errors.go",
			"errorResponse",
			writeErr,
		)
		w.WriteHeader(500)
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error, metadataErr MetadataErr) {
	errorResponse(w, r, http.StatusBadRequest, err.Error(), err, metadataErr)
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request, err error, metadataErr MetadataErr) {
	message := "Forbidden Access"
	if err != nil {
		message = err.Error()
	}
	metadataErr.Message = message
	errorResponse(w, r, http.StatusForbidden, message, err, metadataErr)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request, metadataErr MetadataErr) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	metadataErr.Message = message
	errorResponse(w, r, http.StatusUnauthorized, message, nil, metadataErr)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, metadataErr MetadataErr) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	metadataErr.Message = message
	errorResponse(w, r, http.StatusMethodNotAllowed, message, nil, metadataErr)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, metadataErr MetadataErr) {
	message := "The requested resource could not be found"
	metadataErr.Message = message
	errorResponse(w, r, http.StatusNotFound, message, nil, metadataErr)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string, metadataErr MetadataErr) {
	metadataErr.Message = "validation failed"
	errorResponse(w, r, http.StatusUnprocessableEntity, errors, nil, metadataErr)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error, metadataErr MetadataErr) {
	message := "The server encountered a problem and could not process your request"
	metadataErr.Message = message
	errorResponse(w, r, http.StatusInternalServerError, message, err, metadataErr)
}
