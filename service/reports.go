package service

import (
	"errors"
	"fmt"

	"github.com/heroku/go-getting-started/repository"

	"github.com/heroku/go-getting-started/models"
)

// AddModReport to posts table
func AddModReport(repos *repository.Repository, report models.Report) (int, error) {
	id, err := repos.Reports.CreateReport(report)
	if err != nil {
		return 0, fmt.Errorf("DB can't add report: %w", err)
	}
	return id, nil
}

// UpdateReport from posts
func UpdateReport(repos *repository.Repository, report models.Report) error {
	if err := repos.Reports.UpdateReport(report); err != nil {
		return fmt.Errorf("DB can't Update Report: %w", err)
	}
	return nil
}

// GetReportsByUserId from posts, comments and likes tables
func GetReportsByUserId(repos *repository.Repository, userId int) ([]models.Report, error) {
	post, err := repos.Reports.GetReportsByUserId(userId)
	if err != nil {
		fmt.Println(err.Error())
		return post, errors.New("can't get Reports by userId")

	}
	return post, nil
}

// GetReportsById
func GetReportsById(repos *repository.Repository, reportId int) (models.Report, error) {
	post, err := repos.Reports.GetReportsById(reportId)
	if err != nil {
		fmt.Println(err.Error())
		return post, errors.New("can't get Reports by Id")

	}
	return post, nil
}

// GetAllPosts from posts and likes tables
func GetAllReports(repos *repository.Repository) ([]models.Report, error) {
	allPosts, err := repos.Reports.GetAllReports()
	if err != nil {
		fmt.Println(err.Error())
		return allPosts, errors.New("can't get all Reports")

	}
	return allPosts, nil
}
