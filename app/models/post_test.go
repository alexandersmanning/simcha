package models

import (
	"testing"
)

func clearPosts(t *testing.T) {
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
		rows, err := db.Query("SELECT COUNT(*) FROM posts")

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
			Body:  "Test Body",
			Title: "Test Title",
		}

		err := db.CreatePost(p)

		if err != nil {
			t.Fatal(err)
		}

		postHelper(1, t)
	})
}

func TestCreatePost(t *testing.T) {
	t.Run("Handle posts without associated user", func(t *testing.T) {
		clearPosts(t)
		_, err := db.Query(
			`INSERT INTO posts (body, title) VALUES
			('body_1', 'title_1'),
			('body_2', 'title_2')`)

		if err != nil {
			t.Fatal(err)
		}

		posts, err := db.AllPosts()
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
	})

	t.Run("Handles posts with associated user", func(t *testing.T) {
		clearPosts(t)
		clearUsers(t)

		rows, err := db.Query(`
			INSERT INTO users (email) VALUES ('email@fake.com') RETURNING id
		`)

		var userID int

		if err != nil {
			t.Fatal(err)
		}

		for rows.Next() {
			if err := rows.Scan(&userID); err != nil {
				t.Fatal(err)
			}
		}

		_, err = db.Query(`
			INSERT INTO posts (user_id, title, body) VALUES ($1, $2, $3)
		`, userID, "title_3", "body_3")

		if err != nil {
			t.Fatal(err)
		}

		posts, err := db.AllPosts()
		if err != nil {
			t.Fatal(err)
		}

		t.Run("Returns the correct number of elements", func(t *testing.T) {
			if len(posts) != 1 {
				t.Errorf("Expected %d posts, got %d", 1, len(posts))
			}
		})

		t.Run("Elements have the correct User and Body", func(t *testing.T) {
			post := posts[0]
			if post.Author.Email != "email@fake.com" {
				t.Errorf("Expected user email to be %s, got %s", "email@fake.com", post.Author.Email)
			}

			if post.Body != "body_3" {
				t.Errorf("Expected body to be %s, got %s", "body_3", post.Body)
			}
		})
	})

}
