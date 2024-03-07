package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/models"
	"github.com/heroku/go-getting-started/service"
)

func (h *Handler) memberCommentCreate(c *gin.Context) {
	userId, err := getUserId(c)
	var input models.Comment
	input.CreatedBy = userId
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := service.AddComment(h.repos, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
