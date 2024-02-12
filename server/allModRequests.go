package server

import (
	"net/http"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// getAllModRequests handler - GET only
func (h *Handler) getAllModRequests(w http.ResponseWriter, r *http.Request) {
	// ---negative cases---
	if r.URL.Path != "/requests" {
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
	allModRequests, err := service.GetAllModRequests(h.repos)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tpl.ExecuteTemplate(w, "allModRequests.html", allModRequests); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}
