package middleware

import (
	"fmt"
	"net/http"

	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				metadataErr := utils.Envelope{
					"file": "recoverPanic.go",
					"func": "RecoverPanic",
					"line": 0,
				}
				w.Header().Set("Connection", "close")
				utils.ServerErrorResponse(w, r, fmt.Errorf("%s", err), metadataErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}