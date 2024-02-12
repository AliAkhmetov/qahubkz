package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// modRequestCreate handler - POST only
func (h *Handler) modReportCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/report/create" {
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
	if user.UserType != "moderator" {
		Errors(w, http.StatusBadRequest, "Incorrect user Type")
	}
	postId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	var report models.Report

	report.CreatedBy = user.Id
	report.CreatedAt = time.Now()
	report.PostId = postId
	report.ModeratorMsg = r.FormValue("report-text")
	report.Status = "created"
	_, err = service.AddModReport(h.repos, report)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/my-reports", http.StatusFound)
}

// reportStatusChange handler - POST only
func (h *Handler) reportStatusChange(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/report/change" {
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
	reportId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		Errors(w, http.StatusBadRequest, "Incorrect reportId Id ")
	}

	report, err := service.GetReportsById(h.repos, reportId)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
	}
	report.AdminMsg = r.FormValue("admin-text")

	status := r.FormValue("status")
	if status != "accept" && status != "decline" {
		Errors(w, http.StatusBadRequest, "Incorrect modRequest status")
	}
	report.Status = r.FormValue("status")

	err = service.UpdateReport(h.repos, report)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/reports", http.StatusFound)
}
