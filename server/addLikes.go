package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/models"
	"github.com/heroku/go-getting-started/service"
)

// memberLikeForPost
func (h *Handler) memberLikeForPost(c *gin.Context) {
	userId, err := getUserId(c)
	var input models.LikePost
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.CreatedBy = userId
	id, err := service.AddLikePost(h.repos, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// memberLikeForComment
func (h *Handler) memberLikeForComment(c *gin.Context) {
	userId, err := getUserId(c)
	var input models.LikeComment
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.CreatedBy = userId
	id, err := service.AddLikeComment(h.repos, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
