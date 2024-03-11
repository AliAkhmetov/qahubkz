package server

import (
	"html/template"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/repository"
	//"github.com/swaggo/gin-swagger"
	//
	// "github.com/heroku/go-getting-started/repository"
)

var (
	tpl *template.Template
)

type Handler struct {
	repos *repository.Repository
}

type errorResponse struct {
	Message string `json:"message"`
}

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"https://www.qahub.kz/": true, // Добавьте сюда другой домен
		"https://www.qahub.kz":  true, // Добавьте сюда другой домен

	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверяем, разрешен ли источник
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Language")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// NewHandler create Handler struct with repos parameter
func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{repos: repos}
}
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(CORSMiddleware())

	auth := router.Group("/auth")
	{
		// Auth Handlers
		auth.POST("/registration", h.gestRegistration)
		auth.POST("/login", h.gestLogin)
		//auth.POST("/logout", h.memberLogout)
	}
	unauth := router.Group("/posts")
	{
		// Auth Handlers
		unauth.GET("/", h.getAllPosts)
		unauth.GET("/:id", h.getPostAndComments)

		//auth.POST("/logout", h.memberLogout)
	}
	api := router.Group("/api", h.userIdentity)
	{
		posts := api.Group("/posts")
		{
			posts.GET("/", h.getAllPosts)
			posts.GET("/:id", h.getPostAndComments)
			posts.POST("/", h.memberPostCreate)
			posts.PUT("/:id", h.memberPostUpdate)
			posts.DELETE("/:id", h.memberPostDelete)

			likes := posts.Group(":id/likes")
			{
				likes.POST("/", h.memberLikeForPost)
			}

			comments := posts.Group(":id/comments")
			{
				comments.POST("/", h.memberCommentCreate)
				comments.POST("/like", h.memberLikeForComment)
			}

		}
	}

	return router
}

//TO DO
// // Middleware for identification the user by cookie
// member := h.identification
// //-------- Admin ----------
// //Users
// http.HandleFunc("/users", member(h.getAllUsers))
// //Users type change
// http.HandleFunc("/v1/user/type/change", member(h.updateUserType))
// //Requests
// http.HandleFunc("/requests", member(h.getAllModRequests))
// //Requests status change
// http.HandleFunc("/v1/request-moderator/status/change", member(h.modReqStatusChange))
// //Reports
// http.HandleFunc("/reports", member(h.getAllReports))
// //Reports status change
// http.HandleFunc("/v1/report/change", member(h.reportStatusChange))
// //-------- Moderator ----------
// //Request
// http.HandleFunc("/v1/request-moderator/create", member(h.modRequestCreate))
// //Reports
// http.HandleFunc("/v1/report/create", member(h.modReportCreate))
// http.HandleFunc("/my-reports", member(h.getMyReports))
// //Delete post
// http.HandleFunc("/v1/post/delete", member(h.memberPostDelete))
