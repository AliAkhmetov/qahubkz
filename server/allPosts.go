package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/service"
)

// getAllPosts handler - GET only
func (h *Handler) getAllPosts(c *gin.Context) {
	userId, err := getUserId(c)
	allPosts, err := service.GetAllPosts(h.repos, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, allPosts)
}
