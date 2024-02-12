package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/heroku/go-getting-started/models"
	"github.com/heroku/go-getting-started/service"
)

func (h *Handler) memberPostCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/post/create" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodPost {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var post models.Post
	post.Title = r.FormValue("postTitle")
	post.Content = r.FormValue("postContent")
	post.AuthorName = user.UserName
	post.CreatedBy = user.Id
	var cats []int
	if r.FormValue("1") == "on" {
		cats = append(cats, 1)
	}
	if r.FormValue("2") == "on" {
		cats = append(cats, 2)
	}
	if r.FormValue("3") == "on" {
		cats = append(cats, 3)
	}
	if r.FormValue("4") == "on" {
		cats = append(cats, 4)
	}

	if len(cats) == 0 {
		cats = []int{1, 2, 3, 4}
	}
	_, err := service.AddPost(h.repos, post, cats)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
	// TODO add context of error or created post
	http.Redirect(w, r, "/posts", http.StatusFound)
}

// memberPostDelete
func (h *Handler) memberPostDelete(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/post/delete" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodPost {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		Errors(w, http.StatusBadRequest, "Incorrect modRequest Id ")
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if user.UserType != "admin" && user.UserType != "moderator" {
		Errors(w, http.StatusBadRequest, "Incorrect user Type")
		return
	}
	err = service.DeletePostById(h.repos, id)
	if err != nil {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts"), http.StatusFound)
}
