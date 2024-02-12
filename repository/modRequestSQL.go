package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/heroku/go-getting-started/models"

	_ "github.com/mattn/go-sqlite3"
)

type modRequestSQL struct {
	db *sql.DB
}

// NewPostsSQL create new database struct
func NewModRequestSQL(db *sql.DB) *modRequestSQL {
	return &modRequestSQL{db: db}
}

// CreateModRequest
func (r *modRequestSQL) CreateModRequest(modRequest models.ModRequest) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (user_id, created_at, status) values (?,?,?) RETURNING id`, modRequestsTable)
	row := r.db.QueryRow(query, modRequest.UserId, modRequest.CreatedAt, modRequest.Status)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetAllModRequests
func (r *modRequestSQL) GetAllModRequests() ([]models.ModRequest, error) {
	var modRequests []models.ModRequest

	query := fmt.Sprintf(`

	SELECT * FROM %s;`, modRequestsTable)

	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no modRequest found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get modRequests: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var UpdatedAt sql.NullTime

		var modRequest models.ModRequest
		if err = rows.Scan(
			&modRequest.Id,
			&modRequest.UserId,
			&modRequest.CreatedAt,
			&UpdatedAt,
			&modRequest.Status,
		); err != nil {
			return nil, fmt.Errorf("can't scan modRequest: %w", err)
		}
		modRequest.UpdatedAt = UpdatedAt.Time
		modRequests = append(modRequests, modRequest)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get modRequest: %w", err)
	}
	return modRequests, nil
}

// UpdateModRequest
func (r *modRequestSQL) UpdateModRequest(modRequestId int, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET status=? where id=?`, modRequestsTable)
	if _, err := r.db.Exec(query, status, modRequestId); err != nil {
		return fmt.Errorf("can't update modRequest: %w", err)
	}
	return nil
}

func (r *modRequestSQL) GetModRequestByUserId(userId int) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE user_id=?`, modRequestsTable)
	row := r.db.QueryRow(query, userId)
	var modRequest models.ModRequest
	if err := row.Scan(
		&modRequest.Id,
	); err != nil {
		return 0, fmt.Errorf("can't scan modRequest: %w", err)
	}
	return modRequest.Id, nil
}
