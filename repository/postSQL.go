package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/heroku/go-getting-started/models"

	_ "github.com/mattn/go-sqlite3"
)

type postSQL struct {
	db *sql.DB
}

// NewPostsSQL create new database struct
func NewPostsSQL(db *sql.DB) *postSQL {
	return &postSQL{db: db}
}

// CreatePost
// INSERT INTO posts (created_by, created_at, title, content) values(  1, "2023-05-01 13:35:04.556898354+06:00" , "post about JS", "JavaScript is a scripting or programming language that allows you to implement complex features on web pages — every time a web page does more than just sit there and display static information for you to look at — displaying timely content updates, interactive maps, animated 2D,3D graphics, scrolling video jukeboxes, etc. — you can bet that JavaScript is probably involved. It is the third layer of the layer cake of standard web technologies, two of which (HTML and CSS) we have covered in much more detail in other parts of the Learning Area.");
func (r *postSQL) CreatePost(post models.Post) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (created_by, created_at, title, content, status) values (?,?,?,?,?) RETURNING id`, postsTable)
	row := r.db.QueryRow(query, post.CreatedBy, post.CreatedAt, post.Title, post.Content, "created")
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// AddCategoriesToPost
// INSERT INTO posts_categories (post_id, category_id) values (1, 2);
func (r *postSQL) AddCategoryToPost(postId, catId int) error {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (post_id, category_id) values (?, ?) RETURNING id`, categoriesToPostsTable)
	row := r.db.QueryRow(query, postId, catId)
	if err := row.Scan(&id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SELECT p.*, group_concat(c.name, ", "), (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as Likes, (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as Dislikes, (SELECT pl.id FROM  posts_likes pl WHERE pl.post_id = p.id and type = true and pl.created_by = 1) as like_id FROM posts p LEFT JOIN posts_categories pc ON p.id = pc.post_id LEFT JOIN categories c ON c.id = pc.category_id group by p.id;

// GetAllPosts
// SELECT p.*, group_concat(c.name, ", "), (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as Likes, (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as Dislikes, FROM posts p LEFT JOIN posts_categories pc ON p.id = pc.post_id LEFT JOIN categories c ON c.id = pc.category_id group by p.id;
func (r *postSQL) GetAllPosts(currentUserId int) ([]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf(`
		SELECT p.*, u.username as username, STRING_AGG(c.name, ', ') as categories, 
		(SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as likes,
		(SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as dislikes,
		(SELECT pl.id FROM posts_likes pl WHERE pl.post_id = p.id and type = true and pl.created_by = $1) as my_like_id 
		FROM posts p 
		LEFT JOIN posts_categories pc ON p.id = pc.post_id 
		LEFT JOIN categories c ON c.id = pc.category_id 
		LEFT JOIN users u ON u.id = p.created_by  
		GROUP BY p.id, u.username;
	`)

	rows, err := r.db.Query(query, currentUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no posts found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var myLikeId sql.NullInt32
		var UpdatedAt sql.NullTime

		var post models.Post
		if err = rows.Scan(
			&post.Id,
			&post.CreatedBy,
			&post.CreatedAt,
			&UpdatedAt,
			&post.Title,
			&post.Status,
			&post.Content,
			&post.AuthorName,
			&post.Categories,
			&post.Likes,
			&post.Dislikes,
			&myLikeId,
		); err != nil {
			return nil, fmt.Errorf("can't scan posts: %w", err)
		}
		post.UpdatedAt = UpdatedAt.Time
		post.MyLikeId = int(myLikeId.Int32)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	return posts, nil
}

// GetPostById
// SELECT p.*, group_concat(c.name, ", "), (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as Likes, (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as Dislikes FROM posts p LEFT JOIN posts_categories pc ON p.id = pc.post_id LEFT JOIN categories c ON c.id = pc.category_id WHERE p.id = 7;
func (r *postSQL) GetPostById(userId int) (models.Post, error) {
	var post models.Post
	query := fmt.Sprintf(`
	SELECT p.*, u.username as username, group_concat(c.name, ", ") as categories, 
	(SELECT Count(*) FROM %s  pl WHERE pl.post_id = p.id and type = true) as likes,
	(SELECT Count(*) FROM %s  pl WHERE pl.post_id = p.id and type = false) as dislikes 
	FROM %s p 
	LEFT JOIN %s pc ON p.id = pc.post_id 
	LEFT JOIN %s c ON c.id = pc.category_id 
	LEFT JOIN %s u ON u.id = p.created_by  
	WHERE p.id = ?`, postsLikesTable, postsLikesTable, postsTable, categoriesToPostsTable, categoriesTable, usersTable)

	row := r.db.QueryRow(query, userId)
	var UpdatedAt sql.NullTime

	err := row.Scan(
		&post.Id,
		&post.CreatedBy,
		&post.CreatedAt,
		&UpdatedAt,
		&post.Title,
		&post.Status,
		&post.Content,
		&post.AuthorName,
		&post.Categories,
		&post.Likes,
		&post.Dislikes)
	if err != nil {
		return post, err
	}
	post.UpdatedAt = UpdatedAt.Time

	return post, nil
}

// DeletePostById
func (r *postSQL) DeletePostById(postId int) error {
	query := fmt.Sprintf(`DELETE FROM %s where id=?`, postsTable)
	if _, err := r.db.Exec(query, postId); err != nil {
		return fmt.Errorf("can't set post: %w", err)
	}
	return nil
}

// DeleteCategoriesToPost
func (r *postSQL) DeleteCategoriesToPost(postId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE post_id=?`, categoriesToPostsTable)
	if _, err := r.db.Exec(query, postId); err != nil {
		return fmt.Errorf("can't set post: %w", err)
	}
	return nil
}
