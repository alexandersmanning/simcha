package models

import (
	"github.com/alexandersmanning/simcha/app/models"
)

func (s *StoreSuite) TestPostCreation() {
	p := models.Post{
		Author: "Test Username",
		Body:   "Test Body",
		Title:  "Test Title",
	}

	err := models.CreatePost(p)

	if err != nil {
		s.T().Error(err)
	}

	var count int
	rows, err := s.db.Query("SELECT COUNT(*) FROM posts")
	defer rows.Close()

	if err != nil {
		s.T().Error(err)
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("Incorrect count, expected 1 got %d", count)
	}
}

func (s *StoreSuite) TestGetAllPosts() {
	_, err := s.db.Query(
		"INSERT INTO posts (body, title) VALUES" +
			"('body_1', 'title_1')," +
			"('body_2', 'title_2')")

	if err != nil {
		s.T().Error(err)
	}

	posts, err := models.GetAllPosts()

	if err != nil {
		s.T().Error(err)
	}

	if len(posts) != 2 {
		s.T().Errorf("Incorrect number of records, expecting 2 got %d", len(posts))
	}

	if post := posts[0]; post.Body != "body_1" {
		s.T().Errorf("Incorrect body, expected 'body_1', got %s", post.Body)
	}

	if post := posts[1]; post.Title != "title_2" {
		s.T().Errorf("Incorrect body, expected 'title_2', got %s", post.Body)
	}
}
