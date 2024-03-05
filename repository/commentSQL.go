package repository

import (
	"database/sql"
	"fmt"

	"github.com/heroku/go-getting-started/models"

	_ "github.com/mattn/go-sqlite3"
)

type commentSQL struct {
	db *sql.DB
}

// New create new database.
func NewCommentSQL(db *sql.DB) *commentSQL {
	return &commentSQL{db: db}
}

// INSERT INTO comments (created_by, created_at, post_id , content) values (1, "2023-05-01 13:35:04.556898354+06:00" , 1, "golang top, js for girls");
// CreateComment
func (r *commentSQL) CreateComment(comment models.Comment) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (created_by, created_at, post_id, content, status) values (?,?,?,?,?) RETURNING id`, commentsTable)

	row := r.db.QueryRow(query, comment.CreatedBy, comment.CreatedAt, comment.PostID, comment.Content, "created")
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// SELECT c.* , u.username as username,(SELECT Count (*) FROM comments_likes cl WHERE cl.comment_id = c.id and type = true) as likes,(SELECT Count (*) FROM comments_likes cl WHERE cl.comment_id = c.id and type = false) as dislikes FROM comments c LEFT JOIN users u ON u.id = c.created_by WHERE post_id=7;
// GetCommentsByPostId

func (r *commentSQL) GetCommentsByPostId(postId int) ([]models.Comment, error) {
	var allComments []models.Comment

	query := fmt.Sprintf(`
    SELECT c.*, u.username as username,
        (SELECT Count(*) FROM %s cl WHERE cl.comment_id = c.id and cl.type = true) as likes,
        (SELECT Count(*) FROM %s cl WHERE cl.comment_id = c.id and cl.type = false) as dislikes
    FROM %s c 
    LEFT JOIN %s u ON u.id = c.created_by
    WHERE c.post_id = $1;`, commentsLikesTable, commentsLikesTable, commentsTable, usersTable)

	rows, err := r.db.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("can't get comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		var UpdatedAt sql.NullTime

		if err = rows.Scan(
			&comment.Id,
			&comment.CreatedBy,
			&comment.CreatedAt,
			&UpdatedAt,
			&comment.PostID,
			&comment.Content,
			&comment.Status,
			&comment.AuthorName, // Убедитесь, что это поле существует в модели Comment
			&comment.Likes,
			&comment.Dislikes,
		); err != nil {
			return nil, fmt.Errorf("can't scan comments: %w", err)
		}
		if UpdatedAt.Valid {
			comment.UpdatedAt = UpdatedAt.Time
		} else {
			// Обработка случая, когда UpdatedAt является NULL
			// Можно присвоить zero value для времени или оставить поле пустым
		}

		allComments = append(allComments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return allComments, nil
}
