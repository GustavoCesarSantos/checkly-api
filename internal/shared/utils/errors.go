package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	ErrMissingJSONValue = errors.New("BODY MUST CONTAIN A JSON VALUE")
	ErrMissingOrInvalidLimitQueryParam = errors.New("MISSING OR INVALID LIMIT QUERY PARAM")
	ErrInvalidLimitQueryParam = errors.New("INVALID LAST ID QUERY PARAM")
	ErrEditConflict = errors.New("EDIT CONFLICT")
	ErrRecordNotFound = errors.New("RECORD NOT FOUND")
	ErrFailedCheckUrl = errors.New("FAILED TO CHECK URL")
)

type ErrorEnvelope struct {
	Error string `json:"error" example:"error message"`
}

func logError(r *http.Request, err error, metadataErr Envelope) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)
    slog.Error(err.Error(), "method", method, "url", url, "meta", fmt.Sprintf("%s", metadataErr))
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, metadataErr Envelope) {
	data := Envelope{"error": message}
	err := WriteJSON(w, status, data, nil)
	if err != nil {
        logError(r, err, metadataErr)
		w.WriteHeader(500)
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error, metadataErr Envelope) {
	errorResponse(w, r, http.StatusBadRequest, err.Error(), metadataErr)
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request, err error, metadataErr Envelope) {
	message := "Forbidden Access" 
	if err != nil {
		message = err.Error()
	}
	errorResponse(w, r, http.StatusForbidden, message, metadataErr)

}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request, metadataErr Envelope) {
    w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message, metadataErr)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, metadataErr Envelope) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message, metadataErr)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, metadataErr Envelope) {
	message := "The requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message, metadataErr)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error, metadataErr Envelope) {
	logError(r, err, metadataErr)
	message := "The server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message, metadataErr)
}