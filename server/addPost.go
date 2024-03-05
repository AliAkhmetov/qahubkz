package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/heroku/go-getting-started/models"
	"github.com/heroku/go-getting-started/service"
)

func (h *Handler) memberPostCreate(c *gin.Context) {
	userId, err := getUserId(c)
	var input models.Post
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := service.AddPost(h.repos, input, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
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
