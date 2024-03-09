package models

import (
	"errors"
	"time"
)

type NewUser struct {
	Id       int    `json:"id" db:"id"`
	Email    string `json:"email"  db:"email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id           int       `json:"id" db:"id"`
	Email        string    `json:"email"  db:"email"`
	UserName     string    `json:"userName"  db:"username"`
	PassHash     string    `json:"password_hash"  db:"password_hash"`
	Token        string    `json:"token"  db:"token"`
	UserType     string    `json:"userType"  db:"user_type"` // user, moderator, admin
	ModRequested string    `json:"modRequested"  db:"mod_requested"`
	ExpireAt     time.Time `json:"expireAt"  db:"expire_at"`
}

type Post struct {
	Id            int       `json:"id"  db:"id"`
	CreatedBy     int       `json:"createdBy"  db:"created_by"`
	AuthorName    string    `json:"authorName"  db:"username"`
	CreatedAt     time.Time `json:"createdAt"  db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt"  db:"updated_at"`
	Title         string    `json:"title"  db:"title"`
	Content       string    `json:"content"  db:"content"`
	Categories    string    `json:"categories" db:"categories"`
	CategoriesInt []int     `json:"categoriesInt"`
	Likes         int       `json:"likes"  db:"likes"`
	Dislikes      int       `json:"dislikes"  db:"dislikes"`
	LikedByMe     bool      `json:"likedByMe"  db:"liked_by_me"`
	DislikedByMe  bool      `json:"dislikedByMe"  db:"disliked_by_me"`
	Language      string    `json:"language" db:"language"`

	Status    string `json:"status"  db:"status"`
	ReadTime  int    `json:"readTime"  db:"read_time"`
	ImageLink string `json:"imageLink"  db:"image_link"`
}

type ModRequest struct {
	Id        int       `json:"id"  db:"id"`
	UserId    int       `json:"userId"  db:"user_id"`
	CreatedAt time.Time `json:"createdAt"  db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt"  db:"updated_at"`
	Status    string    `json:"status"  db:"status"`
}

type Report struct {
	Id           int       `json:"id"  db:"id"`
	CreatedBy    int       `json:"createdBy"  db:"created_by"`
	PostId       int       `json:"postId"  db:"post_id"`
	CreatedAt    time.Time `json:"createdAt"  db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt"  db:"updated_at"`
	ModeratorMsg string    `json:"moderatorMsg"  db:"moderator_msg"`
	AdminMsg     string    `json:"adminMsg"  db:"admin_msg"`
	Status       string    `json:"status"  db:"status"`
}

type Comment struct {
	Id           int       `json:"id"  db:"id"`
	CreatedBy    int       `json:"createdBy"  db:"created_by"`
	AuthorName   string    `json:"authorName"  db:"username"`
	CreatedAt    time.Time `json:"createdAt"  db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt"  db:"updated_at"`
	PostID       int       `json:"postId"  db:"post_id"`
	Content      string    `json:"content"  db:"content"`
	Likes        int       `json:"likes"  db:"likes"`
	Dislikes     int       `json:"dislikes"  db:"dislikes"`
	Status       string    `json:"status"  db:"status"`
	LikedByMe    bool      `json:"likedByMe"  db:"liked_by_me"`
	DislikedByMe bool      `json:"dislikedByMe"  db:"disliked_by_me"`
}

type PostAndComments struct {
	Post_info Post      `json:"post_info"`
	Comments  []Comment `json:"comments"`
	IsAuth    bool      `json:"autorized"`
	UserId    int       `json:"userId"`
	UserType  string    `json:"userType"`
}

type LikePost struct {
	Id        int  `json:"id"  db:"id"`
	CreatedBy int  `json:"createdBy"  db:"created_by"`
	PostID    int  `json:"postId"  db:"post_id"`
	Type      bool `json:"type"  db:"type"`
}

type LikeComment struct {
	Id        int  `json:"id"  db:"id"`
	CreatedBy int  `json:"createdBy"  db:"created_by"`
	CommentID int  `json:"commentId"  db:"comment_id"`
	Type      bool `json:"type"  db:"type"`
}

type Categories struct {
	Id   int    `json:"id"  db:"id"`
	Name string `json:"name"  db:"name"`
}

type RegistrationPage struct {
	SuccessMessage string `json:"successMessage"`
	ErrorMessage   string `json:"errorMessage"`
}

var ErrorUnauthorized = errors.New("Unauthorized")
