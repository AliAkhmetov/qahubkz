package service

import (
	"fmt"

	"github.com/heroku/go-getting-started/repository"

	"github.com/heroku/go-getting-started/models"
)

// GetAllComments from comments and likes tables
func AddLikePost(repos *repository.Repository, like models.LikePost) (int, error) {
	id, err := repos.Likes.AddLikePost(like)
	if err != nil {
		return 0, fmt.Errorf("DB can't add like to post: %w", err)
	}
	return id, nil
}

// AddComment to comments table
func AddLikeComment(repos *repository.Repository, like models.LikeComment) (int, error) {
	id, err := repos.Likes.AddLikeComment(like)
	if err != nil {
		return 0, fmt.Errorf("DB can't add like to comment: %w", err)
	}
	return id, nil
}
