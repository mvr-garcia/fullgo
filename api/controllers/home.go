package controllers

import (
	"net/http"

	"github.com/mvr-garcia/fullgo/api/responses"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JsonResponse(w, http.StatusOK, "Welcome to this awesome API")
}
