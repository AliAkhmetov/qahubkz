package repository

import (
	"database/sql"
	"fmt"
	"strconv"

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
	postId, err := strconv.Atoi(like.PostID)
	if err != nil {
		return 0, err
	}
	likeFromDb, err := r.GetLikeByPostUser(postId, like.CreatedBy)
	query := ""
	if likeFromDb.Id != 0 {
		if likeFromDb.Type != like.Type {
			query = fmt.Sprintf(`UPDATE %s SET type = ?  WHERE id = ?`, postsLikesTable)
			if _, err := r.db.Exec(query, like.Type, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't set like type: %w", err)
			}
			return likeFromDb.Id, nil
		} else {
			query = fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, postsLikesTable)
			if _, err := r.db.Exec(query, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't delete like: %w", err)
			}
			return 0, nil
		}
	} else {
		query = fmt.Sprintf(`INSERT INTO %s (created_by, post_id, type) values (?,?,?) RETURNING id`, postsLikesTable)
	}
	row := r.db.QueryRow(query, like.CreatedBy, like.PostID, like.Type)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// AddLikeComment
// INSERT INTO comments_likes (created_by, comment_id, type) values (1,2,true)
func (r *likeSQL) AddLikeComment(like models.LikeComment) (int, error) {
	var id int
	CommentID, err := strconv.Atoi(like.CommentID)
	if err != nil {
		return 0, err
	}
	likeFromDb, err := r.GetLikeByCommentUser(CommentID, like.CreatedBy)
	query := ""
	if likeFromDb.Id != 0 {
		if likeFromDb.Type != like.Type {
			query = fmt.Sprintf(`UPDATE %s SET type = ?  WHERE id = ?`, commentsLikesTable)
			if _, err := r.db.Exec(query, like.Type, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't set like type: %w", err)
			}
			return likeFromDb.Id, nil
		} else {
			query = fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, commentsLikesTable)
			if _, err := r.db.Exec(query, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't delete like: %w", err)
			}
			return 0, nil
		}
	} else {
		query = fmt.Sprintf(`INSERT INTO %s (created_by, comment_id, type) values (?,?,?) RETURNING id`, commentsLikesTable)
	}
	row := r.db.QueryRow(query, like.CreatedBy, like.CommentID, like.Type)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// GetLikeByPostUser
// SELECT * FROM posts_likes WHERE post_id=1 AND created_by=2
func (r *likeSQL) GetLikeByPostUser(postId, userId int) (models.LikePost, error) {
	var like models.LikePost

	query := fmt.Sprintf("SELECT * FROM %s WHERE post_id= ? AND created_by= ?", postsLikesTable)
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

	query := fmt.Sprintf("SELECT * FROM %s WHERE comment_id= ? AND created_by=?", commentsLikesTable)
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
