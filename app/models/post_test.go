package models

import (
	"testing"

	"github.com/alexandersmanning/simcha/app/shared/database"
)

func clearPosts(t *testing.T) {
	db := database.GetStore()
	_, err := db.Query("DELETE FROM posts")

	if err != nil {
		t.Fatal(err)
	}
}

func TestPostCreation(t *testing.T) {
	clearPosts(t)

	postHelper := func(expected int, t *testing.T) {
		t.Helper()
		var count int
		rows, err := database.GetStore().Query("SELECT COUNT(*) FROM posts")

		if err != nil {
			t.Fatal(err)
		}

		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&count)

			if err != nil {
				t.Fatal(err)
			}

			if count != expected {
				t.Errorf("Expected %d records, got %d", expected, count)
			}
		}
	}

	t.Run("No post", func(t *testing.T) {
		postHelper(0, t)
	})

	t.Run("One post", func(t *testing.T) {
		p := Post{
			Author: "Test Username",
			Body:   "Test Body",
			Title:  "Test Title",
		}

		err := CreatePost(p)

		if err != nil {
			t.Fatal(err)
		}

		postHelper(1, t)
	})
}

func TestCreatePost(t *testing.T) {
	clearPosts(t)
	db := database.GetStore()
	_, err := db.Query(
		"INSERT INTO posts (body, title) VALUES" +
			"('body_1', 'title_1')," +
			"('body_2', 'title_2')")

	if err != nil {
		t.Fatal(err)
	}

	posts, err := GetAllPosts()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Returns the correct number of elements", func(t *testing.T) {
		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
	})

	t.Run("Elements container the correct data", func(t *testing.T) {
		if post := posts[0]; post.Body != "body_1" {
			t.Errorf("Incorrect body, expected 'body_1', got %s", post.Body)
		}

		if post := posts[1]; post.Title != "title_2" {
			t.Errorf("Incorrect body, expected 'title_2', got %s", post.Body)
		}
	})
}
