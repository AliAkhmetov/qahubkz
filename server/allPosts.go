package server

import (
	"errors"
	"net/http"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// getAllPosts handler - GET only
func (h *Handler) getAllPosts(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/posts" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	responseBody := struct {
		Posts             []models.Post `json:"posts"`
		IsAuth            bool          `json:"autorized"`
		UserId            int           `json:"userId"`
		UserName          string        `json:"userName"`
		UserType          string        `json:"userType"`
		UserHasModRequest bool          `json:"userHasRequest"`
	}{}
	tokenStatus := true
	token, err := r.Cookie(CookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			tokenStatus = false
		} else {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	var user models.User
	if tokenStatus {
		User, err := service.Identification(h.repos, token.Value)
		if err != nil {
			tokenStatus = false
		} else {
			tokenStatus = true
			user = User
		}

	}
	if !tokenStatus {
		responseBody.IsAuth = false
		responseBody.UserId = 0
		allPosts, err := service.GetAllPosts(h.repos, 0)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		responseBody.Posts = allPosts
	} else {
		responseBody.IsAuth = true
		responseBody.UserId = user.Id
		responseBody.UserName = user.UserName
		responseBody.UserType = user.UserType
		responseBody.UserHasModRequest = service.UserHasModRequest(h.repos, user.Id)
		allPosts, err := service.GetAllPosts(h.repos, user.Id)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		responseBody.Posts = allPosts
	}

	if err := tpl.ExecuteTemplate(w, "allPosts.html", responseBody); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}
