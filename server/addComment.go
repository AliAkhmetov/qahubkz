package server

import (
	"fmt"
	"net/http"

	"github.com/heroku/go-getting-started/models"
	"github.com/heroku/go-getting-started/service"
)

func (h *Handler) memberCommentCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/comment/create" {
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

	var comm models.Comment
	comm.Content = r.FormValue("comment-text")
	comm.PostID = r.FormValue("post-id")
	comm.AuthorName = user.UserName
	comm.CreatedBy = user.Id
	_, err := service.AddComment(h.repos, comm)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	// w.Write([]byte(fmt.Sprintf("%d", id)))
	// w.WriteHeader(http.StatusCreated)

	http.Redirect(w, r, fmt.Sprintf("/post-page?id=%s", comm.PostID), http.StatusFound)
}
