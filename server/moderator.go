package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// modRequestCreate handler - POST only
func (h *Handler) modRequestCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/request-moderator/create" {
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
	if user.UserType != "user" {
		Errors(w, http.StatusBadRequest, "Incorrect user Type")
	}
	var modRequest models.ModRequest
	modRequest.UserId = user.Id
	modRequest.CreatedAt = time.Now()
	modRequest.Status = "created"

	_, err := service.AddModRequest(h.repos, modRequest)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/posts", http.StatusFound)
}

// modReqStatusChange handler - POST only
func (h *Handler) modReqStatusChange(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/request-moderator/status/change" {
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
	modRequestId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		Errors(w, http.StatusBadRequest, "Incorrect modRequest Id ")
	}
	modRequestUserId, err := strconv.Atoi(r.FormValue("userId"))
	if err != nil {
		Errors(w, http.StatusBadRequest, "Incorrect modRequest Id ")
	}
	status := r.FormValue("status")
	if status != "accept" && status != "decline" {
		Errors(w, http.StatusBadRequest, "Incorrect modRequest status")
	}

	err = service.UpdateModRequet(h.repos, modRequestId, modRequestUserId, status)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/posts", http.StatusFound)
}
