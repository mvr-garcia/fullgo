package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mvr-garcia/fullgo/api/auth"
	"github.com/mvr-garcia/fullgo/api/models"
	"github.com/mvr-garcia/fullgo/api/responses"
	"github.com/mvr-garcia/fullgo/api/utils"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.SaveUser(s.DB)
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, formatedError)
		return
	}

	w.Header().Set("location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JsonResponse(w, http.StatusCreated, userCreated)
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	users, err := user.FindAllUsers(s.DB)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	responses.JsonResponse(w, http.StatusOK, users)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	userGotten, err := user.FindUserByID(s.DB, uint32(uid))
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	responses.JsonResponse(w, http.StatusOK, userGotten)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

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

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if tokenID != uint32(uid) {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateAUser(s.DB, uint32(uid))
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		responses.ErrorResponse(w, http.StatusInternalServerError, formatedError)
		return
	}

	responses.JsonResponse(w, http.StatusOK, updatedUser)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = user.DeleteAUser(s.DB, uint32(uid))
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JsonResponse(w, http.StatusNoContent, "")
}
