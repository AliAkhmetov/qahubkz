package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/service"
)

// getAllPosts handler - GET only
func (h *Handler) getAllPosts(c *gin.Context) {
	language := c.GetHeader("language")
	fmt.Println(language)
	userId, err := getUserId(c)
	allPosts, err := service.GetAllPosts(h.repos, userId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, allPosts)
}
