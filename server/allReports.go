package server

import (
	"net/http"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// getAllReports handler - GET only
func (h *Handler) getAllReports(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/reports" {
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
	myReports, err := service.GetAllReports(h.repos)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tpl.ExecuteTemplate(w, "allReports.html", myReports); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// getMyReports handler - GET only
func (h *Handler) getMyReports(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/my-reports" {
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

	if user.UserType != "moderator" {
		Errors(w, http.StatusBadRequest, "Incorrect user Type")
		return
	}
	myReports, err := service.GetReportsByUserId(h.repos, user.Id)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tpl.ExecuteTemplate(w, "myReports.html", myReports); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}
