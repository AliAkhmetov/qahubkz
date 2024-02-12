package server

import (
	"net/http"
	"strconv"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// getAllUsers handler - GET only
func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/users" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if user.UserType != "admin" {
		Errors(w, http.StatusBadRequest, "Incorrect user Type")
		return
	}
	allUsers, err := service.GetAllUsers(h.repos)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tpl.ExecuteTemplate(w, "allUsers.html", allUsers); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// modReqStatusChange handler - POST only
func (h *Handler) updateUserType(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/user/type/change" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodPost {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if user.UserType != "admin" {
		Errors(w, http.StatusBadRequest, "Incorrect user Type")
		return
	}
	userId, err := strconv.Atoi(r.FormValue("user-id"))
	if err != nil {
		Errors(w, http.StatusBadRequest, "Incorrect userId Id ")
		return
	}
	userNewType := r.FormValue("type")

	if userNewType != "user" && userNewType != "moderator" {
		Errors(w, http.StatusBadRequest, "Incorrect user new Type")
	}

	if err := service.UpdateUserType(h.repos, userId, userNewType); err != nil {
		Errors(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}
