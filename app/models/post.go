package models

import (
	"github.com/alexandersmanning/simcha/app/shared/database"
)

type Post struct {
	Author string `json:"author"`
	Body   string `json:"body"`
	Title  string `json:"title"`
}

func GetAllPosts() ([]Post, error) {
	rows, err := database.GetStore().Query("SELECT body, title FROM posts")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		post := Post{}
		if err := rows.Scan(&post.Body, &post.Title); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(p Post) error {
	_, err := database.GetStore().Query(
		"INSERT INTO posts(title, body) VALUES($1, $2) RETURNING id",
		p.Title, p.Body)

	return err
}
