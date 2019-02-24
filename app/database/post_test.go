package database

import (
	"fmt"
	"github.com/alexandersmanning/simcha/app/models"
	"strings"
	"testing"
	"time"
)

func clearPosts(t *testing.T) {
	_, err := db.Query("DELETE FROM posts")

	if err != nil {
		t.Fatal(err)
	}
}

func TestPostCreation(t *testing.T) {
	clearPosts(t)
	clearUsers(t)

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
		u := makeTestUser(t)
		p := models.Post{
			Body:   "Test Body",
			Title:  "Test Title",
			Author: *u,
		}

		err := db.CreatePost(&p)

		if err != nil {
			t.Fatal(err)
		}

		postHelper(1, t)
	})
}

func TestCreatePost(t *testing.T) {
	t.Run("Handle posts without associated user", func(t *testing.T) {
		clearPosts(t)
		clearUsers(t)

		u := makeTestUser(t)

		vals := []interface{}{}
		testData := []struct {
			body     string
			title    string
			id       int
			created  time.Time
			modified time.Time
		}{
			{"body_1", "title_1", u.Id, time.Now().UTC(), time.Now().UTC()},
			{"body_2", "title_2", u.Id, time.Now().UTC(), time.Now().UTC()},
			{"body_3", "title_3", u.Id, time.Now().UTC(), time.Now().UTC()},
		}

		sqlString := `INSERT INTO posts(body, title, user_id, created_at, modified_at) VALUES `
		for idx, row := range testData {
			sqlString += "("
			for i := idx*5 + 1; i <= idx*5+5; i++ {
				sqlString += fmt.Sprintf("$%d, ", i)
			}
			sqlString = strings.TrimSuffix(sqlString, ", ")
			sqlString += "),"
			vals = append(vals, row.body, row.title, row.id, row.created, row.modified)
		}

		sqlString = strings.TrimSuffix(sqlString, ",")
		stmt, _ := db.Prepare(sqlString)
		_, err := stmt.Exec(vals...)
		if err != nil {
			t.Fatal(err)
		}

		posts, err := db.AllPosts()
		if err != nil {
			t.Fatal(err)
		}
		t.Run("Returns the correct number of elements", func(t *testing.T) {
			if len(posts) != 3 {
				t.Errorf("Expected 3 posts, got %d", len(posts))
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

		u := makeTestUser(t)

		_, err := db.Query(`
			INSERT INTO posts (user_id, title, body, created_at, modified_at) VALUES ($1, $2, $3, $4, $5)
		`, u.Id, "title_3", "body_3", time.Now().UTC(), time.Now().UTC())

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
			if post.Author.Email != u.Email {
				t.Errorf("Expected user email to be %s, got %s", u.Email, post.Author.Email)
			}

			if post.Body != "body_3" {
				t.Errorf("Expected body to be %s, got %s", "body_3", post.Body)
			}
		})
	})
}

func TestEditPost(t *testing.T) {
	clearPosts(t)
	clearUsers(t)

	u := makeTestUser(t)
	timestamp := time.Now().UTC()
	rows, err := db.Query(
		`INSERT INTO posts (user_id, title, body, modified_at, created_at) VALUES ($1, $2, $3, $4, $4) RETURNING id`,
		u.Id, "fakeTitle", "fakeBody", timestamp)

	defer rows.Close()

	if err != nil {
		t.Fatal(err)
	}

	var id int
	for rows.Next() {
		err := rows.Scan(&id)

		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("It updates the post", func(t *testing.T) {
		p := models.Post{Id: id, Title: "UpdatedPost", Body: "UpdatedBody" }
		p.SetTimestamps()

		e := db.EditPost(&p)
		if e != nil {
			t.Fatal(e)
		}

		foundP := models.Post{}
		rows, err := db.Query("SELECT title, body, created_at, modified_at FROM posts WHERE id = $1", id)
		if err != nil {
			t.Fatal(err)
		}

		for rows.Next() {
			err := rows.Scan(&foundP.Title, &foundP.Body, &foundP.CreatedAt, &foundP.ModifiedAt)
			if err != nil {
				t.Fatal(err)
			}
		}

		if foundP.Title != p.Title {
			t.Errorf("Expected %s, got %s", p.Title, foundP.Title)
		}

		if foundP.ModifiedAt == foundP.CreatedAt {
			t.Error("Expected modified date to be different than created date")
		}
	})
}

func TestGetPostByID(t *testing.T) {
	u := makeTestUser(t)
	var id string
	insertPost := models.Post{Title: "fakeTitle", Body: "fakeBody", Author: *u}
	insertPost.SetTimestamps()

	rows, err := db.Query(`
	  INSERT INTO posts (user_id, title, body, created_at, modified_at) VALUES ($1, $2, $3, $4, $5) RETURNING id
   `, insertPost.Author.Id, insertPost.Title, insertPost.Body, insertPost.CreatedAt, insertPost.ModifiedAt)

	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
	}

	p, err := db.GetPostById(id)

	if err != nil {
		t.Fatal(err)
	}

	if p.Title != insertPost.Title || p.Body != insertPost.Body || p.Author.Id != insertPost.Author.Id {
		t.Errorf("Expected post title and body to be %s and %s, got %s and %s", insertPost.Title, insertPost.Body, p.Title, p.Body)
	}
}
