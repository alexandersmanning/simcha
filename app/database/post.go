package database

import "github.com/alexandersmanning/simcha/app/models"

//PostStore is the store interface for Posts
type PostStore interface {
	AllPosts() ([]*models.Post, error)
	CreatePost(p models.PostAction) error
	DeletePost(id string) error
	EditPost(p models.PostAction) error
	GetPostById(id string) (*models.Post, error)
}

//AllPosts queries the posts table and returns a slice of Post objects, or and error
func (db *DB) AllPosts() ([]*models.Post, error) {
	rows, err := db.Query(`
		SELECT posts.id,
		       users.id,
		       users.email,
		       posts.body,
		       posts.title,
		       posts.created_at,
		       posts.modified_at
		FROM posts
		LEFT JOIN users ON users.id = posts.user_id
		ORDER BY posts.modified_at DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Author.Id, &post.Author.Email, &post.Body, &post.Title, &post.CreatedAt, &post.ModifiedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

// Returns the Post and Related Author
func (db *DB) GetPostById(id string) (*models.Post, error) {
	var post models.Post
	rows, err := db.Query(`
		SELECT users.id,
		       users.email,
		       posts.title,
		       posts.body,
		       posts.created_at,
		       posts.modified_at
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id = $1
	`,id)
	if err != nil {
		return &post, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&post.Author.Id,
			&post.Author.Email,
			&post.Title,
			&post.Body,
			&post.CreatedAt,
			&post.ModifiedAt,
		); err != nil {
			return &post, err
		}
	}

	return &post, nil
}

//CreatePost creates a new Post object, and returns an ID of the created object
func (db *DB) CreatePost(p models.PostAction) error {
	post := p.Post()
	post.SetTimestamps()

	rows, err := db.Query(
		`INSERT INTO posts(user_id, title, body, created_at, modified_at)
			   VALUES($1, $2, $3, $4, $5)
			   RETURNING id`,
		post.Author.Id, post.Title, post.Body, post.CreatedAt, post.ModifiedAt)

	var id int
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return err
		}
	}

	p.SetID(id)

	return err
}

func (db *DB) EditPost(p models.PostAction) error {
	p.SetTimestamps()
	post := p.Post()

	_, err := db.Query(
		`UPDATE posts SET title = $2, body = $3, modified_at = $4 WHERE id = $1`,
		post.Id, post.Title, post.Body, post.ModifiedAt)

	return err
}

func (db *DB) DeletePost(id string) error {
	_, err := db.Query(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
