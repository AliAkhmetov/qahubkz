package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

// getOnePostAndComments handler - GET only
// Query selectors: id={int}

func (h *Handler) getPostAndComments(c *gin.Context) {
	userId, err := getUserId(c)

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	post, err := service.GetPostById(h.repos, postId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	comments, err := service.GetAllComments(h.repos, postId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	postAndComments := models.PostAndComments{
		Post_info: post,
		Comments:  comments,
	}

	c.JSON(http.StatusOK, postAndComments)
}
