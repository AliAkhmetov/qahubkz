package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/heroku/go-getting-started/models"

	_ "github.com/mattn/go-sqlite3"
)

type usersSQL struct {
	db *sql.DB
}

// NewAuthSQL create new database struct.
func NewAuthSQL(db *sql.DB) *usersSQL {
	return &usersSQL{db: db}
}

// CreateUser in users table | INSERT
func (r *usersSQL) CreateUser(User models.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s ( email, username, password_hash, expire_at, user_type,mod_requested) values (?,?,?,?,?,?) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, User.Email, User.UserName, User.PassHash, time.Now(), User.UserType, false)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *usersSQL) DeleteSuperUser() error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_type=?`, usersTable)
	if _, err := r.db.Exec(query, "admin"); err != nil {
		return fmt.Errorf("can't delete super user: %w", err)
	}
	return nil
}

// GetUser by email from users table | SELECT
func (r *usersSQL) GetUser(Email string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=?", usersTable)
	var token sql.NullString
	err := r.db.QueryRow(query, Email).Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.PassHash,
		&token,
		&user.UserType,
		&user.ModRequested,
		&user.ExpireAt,
	)
	if err != nil {
		return user, err
	}
	user.Token = token.String
	if err != nil {
		return user, fmt.Errorf("can't get user: %w", err)
	}

	return user, nil
}

// AddToken in users table | UPDATE
func (r *usersSQL) AddToken(User models.User) error {
	query := fmt.Sprintf(`UPDATE %s SET token = ?, expire_at = ?  WHERE id = ?`, usersTable)

	if _, err := r.db.Exec(query, User.Token, User.ExpireAt, User.Id); err != nil {
		return fmt.Errorf("can't add token: %w", err)
	}

	return nil
}

// DeleteToken in users table | UPDATE
func (r *usersSQL) DeleteToken(User models.User) error {
	query := fmt.Sprintf(`UPDATE %s SET token = ?, expire_at = ?  WHERE id = ?`, usersTable)

	if _, err := r.db.Exec(query, nil, time.Now(), User.Id); err != nil {
		return fmt.Errorf("can't delete token: %w", err)
	}

	return nil
}

// GetUserByToken from users table  | SELECT
func (r *usersSQL) GetUserByToken(Token string) (models.User, error) {
	var user models.User
	var token sql.NullString
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=?", usersTable)
	err := r.db.QueryRow(query, Token).Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.PassHash,
		&token,
		&user.UserType,
		&user.ModRequested,
		&user.ExpireAt,
	)
	if err != nil {
		return user, models.ErrorUnauthorized
	}
	user.Token = token.String

	if err != nil {
		return user, fmt.Errorf("can't get user by token: %w", err)
	}

	return user, nil
}

func (r *usersSQL) GetAllUsers() ([]models.User, error) {
	var allUsers []models.User
	query := fmt.Sprintf("SELECT * from %s;", usersTable)
	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no users found:%w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		var token sql.NullString
		if err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.UserName,
			&user.PassHash,
			&token,
			&user.UserType,
			&user.ModRequested,
			&user.ExpireAt,
		); err != nil {
			return nil, fmt.Errorf("can't scan users: %w", err)
		}
		user.Token = token.String
		allUsers = append(allUsers, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get users: %w", err)
	}
	return allUsers, nil
}

func (r *usersSQL) UpdateUserType(userId int, userType string) error {
	query := fmt.Sprintf(`UPDATE %s SET user_type = ? WHERE id = ?`, usersTable)

	if _, err := r.db.Exec(query, userType, userId); err != nil {
		return fmt.Errorf("can't update user type: %w", err)
	}

	return nil
}