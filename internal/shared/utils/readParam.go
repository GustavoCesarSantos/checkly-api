package utils

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ReadParam(r *http.Request, paramName string) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	value := params.ByName(paramName)
	if value == "" {
	return "", errors.New("missing parameter: " + paramName)
	}
	return value, nil
}