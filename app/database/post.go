package database

import "github.com/alexandersmanning/simcha/app/models"

//PostStore is the store interface for Posts
type PostStore interface {
	AllPosts() ([]*models.Post, error)
	CreatePost(p models.Post) error
}

//AllPosts queries the posts table and returns a slice of Post objects, or and error
func (db *DB) AllPosts() ([]*models.Post, error) {
	rows, err := db.Query(`
		SELECT COALESCE(users.id, '0'), COALESCE(users.email, 'REMOVED'), posts.body, posts.title
		FROM posts
		LEFT JOIN users on users.id = posts.user_id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Author.Id, &post.Author.Email, &post.Body, &post.Title); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

//CreatePost creates a new Post object, and returns an ID of the created object
func (db *DB) CreatePost(p models.Post) error {
	_, err := db.Query(
		"INSERT INTO posts(user_id, title, body) VALUES($1, $2, $3) RETURNING id",
		p.Author.Id, p.Title, p.Body)

	return err
}
