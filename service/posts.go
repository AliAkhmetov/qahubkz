package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/heroku/go-getting-started/repository"

	"github.com/heroku/go-getting-started/models"
)

// GetAllPosts from posts and likes tables
func GetAllPosts(repos *repository.Repository, currentUserId int) ([]models.Post, error) {
	allPosts, err := repos.Posts.GetAllPosts(currentUserId)
	if err != nil {
		fmt.Println(err.Error())
		return allPosts, errors.New("can't get all posts")

	}
	return allPosts, nil
}

// GetPostById from posts, comments and likes tables
func GetPostById(repos *repository.Repository, postId, userId int) (models.Post, error) {
	post, err := repos.Posts.GetPostById(postId, userId)
	if err != nil {
		fmt.Println(err.Error())
		return post, errors.New("can't get post by id")

	}
	return post, nil
}

// AddPost to posts table
func AddPost(repos *repository.Repository, post models.Post, userId int) (int, error) {
	post.CreatedAt = time.Now()
	post.CreatedBy = userId
	id, err := repos.Posts.CreatePost(post)
	if err != nil {
		return 0, fmt.Errorf("DB can't add post: %w", err)
	}
	for _, catId := range post.CategoriesInt {
		if err := repos.Posts.AddCategoryToPost(id, catId); err != nil {
			return 0, fmt.Errorf("DB can't add category: %w", err)
		}
	}
	return id, nil
}

// DeletePostById from posts
func DeletePostById(repos *repository.Repository, id int) error {
	if err := repos.Posts.DeleteCategoriesToPost(id); err != nil {
		return fmt.Errorf("DB can't delete Categories: %w", err)
	}
	err := repos.Posts.DeletePostById(id)
	if err != nil {
		return errors.New("can't  delete post")
	}
	return nil
}
