package models

//Post is a struct for creating a simple blog post
type Post struct {
	Author User   `json:"author"`
	Body   string `json:"body"`
	Title  string `json:"title"`
}

//PostStore is the store interface for Posts
type PostStore interface {
	AllPosts() ([]*Post, error)
	CreatePost(p Post) error
}

//AllPosts queries the posts table and returns a slice of Post objects, or and error
func (db *DB) AllPosts() ([]*Post, error) {
	rows, err := db.Query(`
		SELECT COALESCE(users.id, '0'), COALESCE(users.email, 'REMOVED'), posts.body, posts.title
		FROM posts
		LEFT JOIN users on users.id = posts.user_id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		post := Post{}
		if err := rows.Scan(&post.Author.ID, &post.Author.Email, &post.Body, &post.Title); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

//CreatePost creates a new Post object, and returns an ID of the created object
func (db *DB) CreatePost(p Post) error {
	_, err := db.Query(
		"INSERT INTO posts(user_id, title, body) VALUES($1, $2, $3) RETURNING id",
		p.Author.ID, p.Title, p.Body)

	return err
}
