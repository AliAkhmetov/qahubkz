package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/heroku/go-getting-started/service"

	"github.com/heroku/go-getting-started/models"
)

const CookieName = "token"

func (h *Handler) gestRegistration(c *gin.Context) {
	var input models.NewUser

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := service.Registration(h.repos, input.Username, input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) gestLogin(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		fmt.Println("qweqwe")

		return
	}
	fmt.Println(input)

	user, err := service.Authentication(h.repos, input.Email, input.Password)
	res := struct {
		Value    string    `json:"token"`
		UserId   int       `json:"userID"`
		Expires  time.Time `json:"expires"`
		ErrorMsg string    `json:"error"`
	}{
		Value:   "",
		Expires: user.ExpireAt,
	}
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := service.Authorization(h.repos, user)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	res.Value = token
	res.UserId = user.Id

	c.JSON(http.StatusOK, res)
}

// memberLogout handler -GET only
func (h *Handler) memberLogout(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := service.Logout(h.repos, user); err != nil {
		fmt.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusFound)
}
