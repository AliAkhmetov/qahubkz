package service

import (
	"errors"
	"fmt"

	"github.com/heroku/go-getting-started/repository"

	"github.com/heroku/go-getting-started/models"
)

// GetAllUsers from posts and likes tables
func GetAllUsers(repos *repository.Repository) ([]models.User, error) {
	allUsers, err := repos.Authorization.GetAllUsers()
	if err != nil {
		fmt.Println(err.Error())
		return allUsers, errors.New("can't get all posts")

	}
	return allUsers, nil
}

func UpdateUserType(repos *repository.Repository, userId int, UserType string) error {
	err := repos.Authorization.UpdateUserType(userId, UserType)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("can't get all posts")
	}
	return nil
}
