package server

import (
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

// memberPostUpdate
func (h *Handler) memberPostUpdate(c *gin.Context) {
	userId, err := getUserId(c)

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var input models.Post
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.Id = postId

	if err := service.UpdatePost(h.repos, input, userId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": postId,
	})
}

// memberPostDelete
func (h *Handler) memberPostDelete(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = service.DeletePostById(h.repos, postId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": 0,
	})
}
