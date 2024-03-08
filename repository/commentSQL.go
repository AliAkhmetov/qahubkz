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
	query := fmt.Sprintf(`INSERT INTO %s (created_by, created_at, post_id, content, status) values ($1,$2,$3,$4,$5) RETURNING id`, commentsTable)

	row := r.db.QueryRow(query, comment.CreatedBy, comment.CreatedAt, comment.PostID, comment.Content, "created")
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// SELECT c.* , u.username as username,(SELECT Count (*) FROM comments_likes cl WHERE cl.comment_id = c.id and type = true) as likes,(SELECT Count (*) FROM comments_likes cl WHERE cl.comment_id = c.id and type = false) as dislikes FROM comments c LEFT JOIN users u ON u.id = c.created_by WHERE post_id=7;
// GetCommentsByPostId

func (r *commentSQL) GetCommentsByPostId(postId, userId int) ([]models.Comment, error) {
	var allComments []models.Comment

	query := fmt.Sprintf(`
    SELECT c.*, u.username as username,
        (SELECT Count(*) FROM %s cl WHERE cl.comment_id = c.id and cl.type = true) as likes,
        (SELECT Count(*) FROM %s cl WHERE cl.comment_id = c.id and cl.type = false) as dislikes,
		(SELECT cl.id FROM comments_likes cl WHERE  cl.comment_id = c.id and type = true and cl.created_by = $1) as liked_by_me,
		(SELECT cl.id FROM comments_likes cl WHERE  cl.comment_id = c.id and type = false and cl.created_by = $2) as disliked_by_me 
    FROM %s c 
    LEFT JOIN %s u ON u.id = c.created_by
    WHERE c.post_id = $3;`, commentsLikesTable, commentsLikesTable, commentsTable, usersTable)

	rows, err := r.db.Query(query, userId, userId, postId)
	if err != nil {
		return nil, fmt.Errorf("can't get comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		var UpdatedAt sql.NullTime
		var LikedByMe sql.NullInt32
		var DislikedByMe sql.NullInt32

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
			&LikedByMe,
			&DislikedByMe,
		); err != nil {
			return nil, fmt.Errorf("can't scan comments: %w", err)
		}
		if UpdatedAt.Valid {
			comment.UpdatedAt = UpdatedAt.Time
		} else {
			// Обработка случая, когда UpdatedAt является NULL
			// Можно присвоить zero value для времени или оставить поле пустым
		}
		if LikedByMe.Int32 >= 1 {
			comment.LikedByMe = true
		}
		if DislikedByMe.Int32 >= 1 {
			comment.DislikedByMe = true
		}
		allComments = append(allComments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return allComments, nil
}
