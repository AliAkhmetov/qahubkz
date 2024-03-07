package repository

import (
	"database/sql"
	"fmt"

	"github.com/heroku/go-getting-started/models"

	_ "github.com/mattn/go-sqlite3"
)

type likeSQL struct {
	db *sql.DB
}

// NewlikeSQL create new database struct
func NewlikeSQL(db *sql.DB) *likeSQL {
	return &likeSQL{db: db}
}

// AddLikePost
// INSERT INTO posts_likes (created_by, post_id, type) values (1,2,true)
func (r *likeSQL) AddLikePost(like models.LikePost) (int, error) {
	var id int

	likeFromDb, _ := r.GetLikeByPostUser(like.PostID, like.CreatedBy)

	query := ""
	if likeFromDb.Id != 0 {
		if likeFromDb.Type != like.Type {
			query = fmt.Sprintf(`UPDATE %s SET type = $1  WHERE id = $2`, postsLikesTable)
			if _, err := r.db.Exec(query, like.Type, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't set like type: %w", err)
			}
			return likeFromDb.Id, nil
		} else {
			query = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, postsLikesTable)
			if _, err := r.db.Exec(query, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't delete like: %w", err)
			}
			return 0, nil
		}
	} else {
		query = fmt.Sprintf(`INSERT INTO %s (created_by, post_id, type) values ($1,$2,$3) RETURNING id`, postsLikesTable)
	}
	row := r.db.QueryRow(query, like.CreatedBy, like.PostID, like.Type)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *likeSQL) AddLikeComment(like models.LikeComment) (int, error) {
	var id int

	// Check if a like already exists for the given comment and user
	existingLike, err := r.GetLikeByCommentUser(like.CommentID, like.CreatedBy)
	if err != nil {
		return 0, err
	}
	if existingLike.Id != 0 {
		if existingLike.Type != like.Type {
			// Update the like type if it has changed
			query := fmt.Sprintf(`UPDATE %s SET type = $1 WHERE id = $2 RETURNING id`, commentsLikesTable)
			err := r.db.QueryRow(query, like.Type, existingLike.Id).Scan(&id)
			if err != nil {
				return 0, fmt.Errorf("can't set like type: %w", err)
			}
		} else {
			// Delete the like if the type is the same (toggle like)
			query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, commentsLikesTable)
			_, err := r.db.Exec(query, existingLike.Id)
			if err != nil {
				return 0, fmt.Errorf("can't delete like: %w", err)
			}
			// Returning 0 to indicate the like was removed
			return 0, nil
		}
	} else {
		// Insert a new like if it doesn't exist
		query := fmt.Sprintf(`INSERT INTO %s (created_by, comment_id, type) VALUES ($1, $2, $3) RETURNING id`, commentsLikesTable)
		err := r.db.QueryRow(query, like.CreatedBy, like.CommentID, like.Type).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("can't insert like: %w", err)
		}
	}

	// Return the id of the like that was affected
	return id, nil
}

// GetLikeByPostUser
// SELECT * FROM posts_likes WHERE post_id=1 AND created_by=2
func (r *likeSQL) GetLikeByPostUser(postId, userId int) (models.LikePost, error) {
	var like models.LikePost

	query := fmt.Sprintf("SELECT * FROM %s WHERE post_id=$1 AND created_by=$2", postsLikesTable)
	err := r.db.QueryRow(query, postId, userId).Scan(
		&like.Id,
		&like.CreatedBy,
		&like.PostID,
		&like.Type,
	)
	if err != nil {
		return like, fmt.Errorf("can't get the like of this post: %w", err)
	}
	return like, nil
}

// GetLikeByCommentUser
// SELECT * FROM comments_likes WHERE post_id=1 AND  created_by=2
func (r *likeSQL) GetLikeByCommentUser(commentId, userId int) (models.LikeComment, error) {
	var like models.LikeComment

	query := fmt.Sprintf("SELECT * FROM %s WHERE comment_id=$1 AND created_by=$2", commentsLikesTable)
	err := r.db.QueryRow(query, commentId, userId).Scan(
		&like.Id,
		&like.CreatedBy,
		&like.CommentID,
		&like.Type,
	)
	if err != nil {
		return like, fmt.Errorf("can't get all the like of this comment: %w", err)
	}
	return like, nil
}
