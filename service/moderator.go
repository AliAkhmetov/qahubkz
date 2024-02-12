package service

import (
	"errors"
	"fmt"

	"github.com/heroku/go-getting-started/repository"

	"github.com/heroku/go-getting-started/models"
)

// GetAllModRequests from mod_requests
func GetAllModRequests(repos *repository.Repository) ([]models.ModRequest, error) {
	allModRequests, err := repos.ModRequests.GetAllModRequests()
	if err != nil {
		fmt.Println(err.Error())
		return allModRequests, errors.New("can't get all ModRequests")

	}
	return allModRequests, nil
}

// AddRequestModerator to mod_requests table
func AddModRequest(repos *repository.Repository, modRequest models.ModRequest) (int, error) {
	id, err := repos.ModRequests.CreateModRequest(modRequest)
	if err != nil {
		return 0, fmt.Errorf("DB can't add modRequest: %w", err)
	}

	return id, nil
}

func UpdateModRequet(repos *repository.Repository, modRequestId, modRequestUserId int, status string) error {
	if status == "accept" {
		err := UpdateUserType(repos, modRequestUserId, "moderator")
		if err != nil {
			return fmt.Errorf("DB can't update user type: %w", err)
		}
	}

	err := repos.ModRequests.UpdateModRequest(modRequestId, status)
	if err != nil {
		return fmt.Errorf("DB can't Update modRequest: %w", err)
	}

	return nil
}

func UserHasModRequest(repos *repository.Repository, userId int) bool {
	id, err := repos.ModRequests.GetModRequestByUserId(userId)
	if err != nil || id == 0 {
		return false
	}
	return true
}
