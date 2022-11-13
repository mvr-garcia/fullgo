package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mvr-garcia/fullgo/api/auth"
	"github.com/mvr-garcia/fullgo/api/models"
	"github.com/mvr-garcia/fullgo/api/responses"
	"github.com/mvr-garcia/fullgo/api/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) SignIn(email, password string) (string, error) {

	user := models.User{}
	err := s.DB.Model(&models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateToken(user.ID)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := s.SignIn(user.Email, user.Password)
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, formatedError)
		return
	}

	responses.JsonResponse(w, http.StatusOK, token)
}
