package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/service"
)

// getAllPosts handler - GET only
func (h *Handler) getAllPosts(c *gin.Context) {

	allPosts, err := service.GetAllPosts(h.repos, 0)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, allPosts)
}
