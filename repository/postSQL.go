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
	query := fmt.Sprintf(`INSERT INTO %s (created_by, created_at, title, content, image_link, read_time, status) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id`, postsTable)
	row := r.db.QueryRow(query, post.CreatedBy, post.CreatedAt, post.Title, post.Content, post.ImageLink, post.ReadTime, "created")
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// AddCategoriesToPost
// INSERT INTO posts_categories (post_id, category_id) values (1, 2);
func (r *postSQL) AddCategoryToPost(postId, catId int) error {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (post_id, category_id) values ($1, $2) RETURNING id`, categoriesToPostsTable)
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
func (r *postSQL) GetAllPosts(currentUserId int, language string) ([]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf(`
		SELECT p.id, p.read_time, p.image_link, p.created_by, p.created_at, p.updated_at, p.title, p.status, u.username as username, STRING_AGG(c.name, ', ') as categories, 
		(SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as likes,
		(SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as dislikes,
		(SELECT pl.id FROM posts_likes pl WHERE pl.post_id = p.id and type = true and pl.created_by = $1) as liked_by_me,
		(SELECT pl.id FROM posts_likes pl WHERE pl.post_id = p.id and type = false and pl.created_by = $2) as disliked_by_me 

		FROM posts p 
		LEFT JOIN posts_categories pc ON p.id = pc.post_id 
		LEFT JOIN categories c ON c.id = pc.category_id 
		LEFT JOIN users u ON u.id = p.created_by
		WHERE p.language = $3
		GROUP BY p.id, u.username
		ORDER BY p.created_at;
	`)

	rows, err := r.db.Query(query, currentUserId, currentUserId, language)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no posts found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var LikedByMe sql.NullInt32
		var DislikedByMe sql.NullInt32
		var UpdatedAt sql.NullTime
		var Categories sql.NullString
		var ImageLink sql.NullString

		var post models.Post
		if err = rows.Scan(
			&post.Id,
			&post.ReadTime,
			&ImageLink,
			&post.CreatedBy,
			&post.CreatedAt,
			&UpdatedAt,
			&post.Title,
			&post.Status,
			&post.AuthorName,
			&Categories,
			&post.Likes,
			&post.Dislikes,
			&LikedByMe,
			&DislikedByMe,
		); err != nil {
			return nil, fmt.Errorf("can't scan posts: %w", err)
		}
		post.Categories = Categories.String
		post.ImageLink = ImageLink.String
		post.UpdatedAt = UpdatedAt.Time
		if LikedByMe.Int32 >= 1 {
			post.LikedByMe = true
		}
		if DislikedByMe.Int32 >= 1 {
			post.DislikedByMe = true
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	return posts, nil
}

// GetPostById

func (r *postSQL) GetPostById(postId, userId int) (models.Post, error) {
	fmt.Println(postId, userId)

	var post models.Post
	query := fmt.Sprintf(`
	SELECT p.id, p.read_time, p.image_link, p.created_by, p.created_at, p.updated_at, p.title, p.status, p.content, u.username as username, 
	STRING_AGG(c.name, ', ') as categories, 
	(SELECT Count(*) FROM %s pl WHERE pl.post_id = p.id and pl.type = true) as likes,
	(SELECT Count(*) FROM %s pl WHERE pl.post_id = p.id and pl.type = false) as dislikes,
	(SELECT pl.id FROM posts_likes pl WHERE pl.post_id = $1 and type = true and pl.created_by = $2) as liked_by_me,
	(SELECT pl.id FROM posts_likes pl WHERE pl.post_id =$3 and type = false and pl.created_by = $4) as disliked_by_me 
	FROM %s p 
	LEFT JOIN %s pc ON p.id = pc.post_id 
	LEFT JOIN %s c ON c.id = pc.category_id 
	LEFT JOIN %s u ON u.id = p.created_by  
	WHERE p.id = $5
	GROUP BY p.id, u.username`, postsLikesTable, postsLikesTable, postsTable, categoriesToPostsTable, categoriesTable, usersTable)

	row := r.db.QueryRow(query, postId, userId, postId, userId, postId)
	//fmt.Printf(row)
	var UpdatedAt sql.NullTime
	var ImageLink sql.NullString
	var Categories sql.NullString
	var LikedByMe sql.NullInt32
	var DislikedByMe sql.NullInt32
	err := row.Scan(
		&post.Id,
		&post.ReadTime,
		&ImageLink,
		&post.CreatedBy,
		&post.CreatedAt,
		&UpdatedAt,
		&post.Title,
		&post.Status,
		&post.Content,
		&post.AuthorName, // Убедитесь, что это поле существует в структуре models.Post
		&Categories,
		&post.Likes,
		&post.Dislikes,
		&LikedByMe,
		&DislikedByMe,
	)
	fmt.Println(LikedByMe.Int32)
	if err != nil {
		return post, err
	}
	if UpdatedAt.Valid {
		post.UpdatedAt = UpdatedAt.Time
	} else {
		// Обработка случая, когда updated_at NULL
	}
	post.ImageLink = ImageLink.String
	post.Categories = Categories.String
	if LikedByMe.Int32 >= 1 {
		post.LikedByMe = true
	}
	if DislikedByMe.Int32 >= 1 {
		post.DislikedByMe = true
	}
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
