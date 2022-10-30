package middlewares

import (
	"errors"
	"net/http"

	"github.com/mvr-garcia/fullgo/api/auth"
	"github.com/mvr-garcia/fullgo/api/responses"
)

func SetMiddlewareJson(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

func SetMiddlewareAuthentication(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidateToken(r)
		if err != nil {
			responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		h(w, r)
	}
}
