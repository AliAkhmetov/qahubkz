package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/heroku/go-getting-started/repository"

	"github.com/heroku/go-getting-started/models"

	"golang.org/x/crypto/bcrypt"
)

// Identification is the act of indicating a person or thing’s identity.
// Authentication is the act of proving  the identity of a computer system user (comparing password entered with the password in DB)
// Authorization is the function of specifying access rights/privileges to resources.

const expireAtTime = 24 * time.Hour

func Registration(repos *repository.Repository, userName, email, password string) (int, error) {
	fmt.Println(userName, email, password)

	if err := checkCreds(userName, email, password); err != nil {
		return http.StatusBadRequest, err
	}
	passHash, err := generatePassHash(password)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, errors.New("Password Hash Error")
	}
	newUser := models.User{
		UserName: userName,
		Email:    email,
		PassHash: passHash,
		UserType: "user",
	}
	id, err := repos.Authorization.CreateUser(newUser)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, errors.New("Unable to save to database")
	}
	return id, nil
}

// Authentication - compare password entered with the password-hash in DB)
func Authentication(repos *repository.Repository, email, password string) (models.User, error) {
	errorUser := models.User{}

	// checking email and pass
	if errEmail := checkEmail(email); errEmail != nil {
		return errorUser, errors.New("Email not valid")
	}
	if errPwd := checkString(password); errPwd != nil {
		return errorUser, errors.New("Pass not valid")
	}

	user, err := repos.Authorization.GetUser(email)
	if err != nil {
		fmt.Println(err.Error())
		return errorUser, errors.New("Incorrect email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		return errorUser, errors.New("Incorrect email or password")
	}

	return user, nil
}

// Authorization - add session token/expiration time to user in DB
func Authorization(repos *repository.Repository, user models.User) (string, error) {
	return AddSessionToken(repos, user)
}

// Logout - delete session token/expiration time to user in DB
func Logout(repos *repository.Repository, user models.User) error {
	err := DeleteSessionToken(repos, user)
	return err
}

func DeleteSessionToken(repos *repository.Repository, user models.User) error {
	user.Token = ""
	user.ExpireAt = time.Now()

	if err := repos.Authorization.DeleteToken(user); err != nil {
		return fmt.Errorf("DB can't delete token: %w", err)
	}
	return nil
}

func AddSessionToken(repos *repository.Repository, user models.User) (string, error) {
	token, err := newToken()
	if err != nil {
		return "", fmt.Errorf("newToken error: %w", err)
	}
	user.Token = token
	user.ExpireAt = time.Now().Add(expireAtTime)

	if err := repos.Authorization.AddToken(user); err != nil {
		return "", fmt.Errorf("DB can't add token: %w", err)
	}
	return token, nil
}

func Identification(repos *repository.Repository, token string) (models.User, error) {
	user, err := GetUserByToken(repos, token)
	if err != nil {
		fmt.Println("gggg")

		return user, err
	}
	return user, nil
}

func GetUserByToken(repos *repository.Repository, token string) (models.User, error) {
	fmt.Println("GetUserByToken")

	user, err := repos.Authorization.GetUserByToken(token)
	if err != nil {
		fmt.Println("DB can't get user by token")

		return user, fmt.Errorf("DB can't get user by token: %w", err)
	}
	if user.ExpireAt.Before((time.Now())) {
		fmt.Println("ExpireAt")

		return user, models.ErrorUnauthorized
	}
	return user, nil
}
