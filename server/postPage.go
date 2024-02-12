package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// getOnePostAndComments handler - GET only
// Query selectors: id={int}
func (h *Handler) getPostAndComments(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/post-page" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	//---positive cases---
	postId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || postId == 0 {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post, err := service.GetPostById(h.repos, postId)
	if err != nil {
		Errors(w, http.StatusNotFound, err.Error())
		return
	}

	comments, err := service.GetAllComments(h.repos, postId)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	postAndComments := models.PostAndComments{
		Post_info: post,
		Comments:  comments,
	}

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
		postAndComments.IsAuth = false
		postAndComments.UserId = 0
	} else {
		postAndComments.IsAuth = true
		postAndComments.UserId = user.Id
		postAndComments.UserType = user.UserType
	}

	if err := tpl.ExecuteTemplate(w, "postPage.html", postAndComments); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}
