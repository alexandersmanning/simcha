package database

import (
	"fmt"
	"github.com/alexandersmanning/simcha/app/models"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setupModels()
	exitCode := m.Run()

	teardownModels()
	os.Exit(exitCode)
}

var db *DB

func setupModels() {
	connString := "dbname=simcha_test sslmode=disable"
	var err error
	db, err = InitDB(connString)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createTestUser(u *models.User, t *testing.T) int {
	var id int

	if u.CreatedAt.IsZero() {
		u.CreatedAt, u.ModifiedAt = time.Now().UTC(), time.Now().UTC()
	}

	rows, err := db.Query(`
		INSERT INTO users(email, password_digest, created_at, modified_at) values ($1, $2, $3, $4)
		RETURNING id
	`, u.Email, u.PasswordDigest, u.CreatedAt, u.ModifiedAt)

	defer rows.Close()

	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
	}

	return id
}

func makeTestUser(t *testing.T) *models.User {
	u := models.User{Email: "fakeuser" + string(rand.Int()), PasswordDigest: "fakedigest"}
	u.Id = createTestUser(&u, t)

	return &u
}

func teardownModels() {
	db.Close()
}
