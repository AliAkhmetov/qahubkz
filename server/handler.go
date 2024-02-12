package server

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/repository"
)

var (
	tpl *template.Template

	port = ":" + os.Getenv("PORT")
)

type Handler struct {
	repos *repository.Repository
}

// NewHandler create Handler struct with repos parameter
func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{repos: repos}
}

// Server func - all handlers
func Server(h *Handler) {
	// Middleware for identification the user by cookie
	member := h.identification

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	// Auth Handlers
	http.HandleFunc("/login", h.gestLogin)
	http.HandleFunc("/registration", h.gestRegistration)
	http.HandleFunc("/logout", member(h.memberLogout))

	// homepage
	http.HandleFunc("/", h.homePage)
	// get all posts
	http.HandleFunc("/posts", h.getAllPosts)
	// gel one post page
	http.HandleFunc("/post-page", h.getPostAndComments)

	//-------- User ----------

	// add post
	http.HandleFunc("/v1/post/create", member(h.memberPostCreate))
	// add comment
	http.HandleFunc("/v1/comment/create", member(h.memberCommentCreate))
	//  add likes
	http.HandleFunc("/v1/post/like", member(h.memberLikeForPost))
	http.HandleFunc("/v1/comment/like", member(h.memberLikeForComment))

	//-------- Admin ----------
	//Users
	http.HandleFunc("/users", member(h.getAllUsers))
	//Users type change
	http.HandleFunc("/v1/user/type/change", member(h.updateUserType))

	//Requests
	http.HandleFunc("/requests", member(h.getAllModRequests))
	//Requests status change
	http.HandleFunc("/v1/request-moderator/status/change", member(h.modReqStatusChange))

	//Reports
	http.HandleFunc("/reports", member(h.getAllReports))
	//Reports status change
	http.HandleFunc("/v1/report/change", member(h.reportStatusChange))

	//-------- Moderator ----------
	//Request
	http.HandleFunc("/v1/request-moderator/create", member(h.modRequestCreate))

	//Reports
	http.HandleFunc("/v1/report/create", member(h.modReportCreate))
	http.HandleFunc("/my-reports", member(h.getMyReports))
	//Delete post
	http.HandleFunc("/v1/post/delete", member(h.memberPostDelete))

	// handle css
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/imgs/", http.StripPrefix("/imgs/", http.FileServer(http.Dir("imgs"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	//
	if port == "" {
		//port = ":8081"
		log.Fatal("$PORT must be set")
	}
	log.Printf("Starting a web server on http://localhost%s/ ", port)

	//port := os.Getenv("PORT")

	// router := gin.New()
	// router.Use(gin.Logger())
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	// router.Static("/static", "static")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	// router.Run(":" + port)
	//port = ":8081"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
