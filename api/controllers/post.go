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

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if uid != post.AuthorID {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	postCreated, err := post.SavePost(s.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ErrorResponse(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JsonResponse(w, http.StatusCreated, postCreated)
}

func (s *Server) GetPosts(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}
	posts, err := post.FindAllPosts(s.DB)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	responses.JsonResponse(w, http.StatusOK, posts)
}

func (s *Server) GetPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	post := models.Post{}
	postReceived, err := post.FindPostByID(s.DB, pid)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	responses.JsonResponse(w, http.StatusOK, postReceived)
}

func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post := models.Post{}
	err = s.DB.Model(&models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ErrorResponse(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	if uid != post.AuthorID {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != postUpdate.AuthorID {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	postUpdate.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID

	postUpdated, err := postUpdate.UpdateAPost(s.DB)
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		responses.ErrorResponse(w, http.StatusInternalServerError, formatedError)
		return
	}

	responses.JsonResponse(w, http.StatusOK, postUpdated)
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post := models.Post{}
	err = s.DB.Model(&models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ErrorResponse(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		responses.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	_, err = post.DeleteAPost(s.DB, pid, uid)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JsonResponse(w, http.StatusOK, "")
}
