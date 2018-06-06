package database

import (
	"testing"
	"github.com/alexandersmanning/simcha/app/models"
	"reflect"
)

func createTestSession(userId int, t *testing.T) models.UserSession {
	token, err := CreateSessionToken()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := db.Query(`
		INSERT INTO user_sessions (user_id, token)
		VALUES ($1, $2) RETURNING ID
	`, userId, token)

	defer rows.Close()

	if err != nil {
		t.Fatal(err)
	}

	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
	}

	return models.UserSession{Id: id, User: models.User{Id: userId}, SessionToken: token }
}

func TestCreateUserSession(t *testing.T) {
	clearUsers(t)

	u := models.User{Email: "email@fake.com", PasswordDigest: "testDigest"}
	u.Id = createTestUser(&u, t)

	us, err := db.CreateUserSession(&u)

	if err != nil {
		t.Fatal(err)
	}

	if us.Id == 0 {
		t.Errorf("Expected a usersession ID, got %d", us.Id)
	}

	if !reflect.DeepEqual(us.User, u) {
		t.Errorf("Expected the user session user to be %v, got %v", u, us.User)
	}

	if us.SessionToken == "" {
		t.Errorf("Expected a token for the user session, got nothing")
	}
}

func TestGetUserBySessionToken(t *testing.T) {
	clearUsers(t)

	uOne := models.User{Email: "email1@fake.com"}
	id := createTestUser(&uOne, t)
	uOne.Id = id
	usOne := createTestSession(id, t)

	t.Run("It returns the user when one user entry exists", func(t *testing.T) {
		user, err := db.GetUserBySessionToken(usOne.User.Id, usOne.SessionToken)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(user, uOne) {
			t.Errorf("Expected %v to equal %v", user, uOne)
		}
	})

	t.Run("It returns an empty user if nothing can be found", func(t *testing.T) {
		user, err := db.GetUserBySessionToken(usOne.User.Id, "non-existent-token")
		if err != nil {
			t.Fatal(err)
		}

		if user.Id != 0 || user.Email != "" {
			t.Errorf("Expected no user to be found, a user with id %d and email %s was returned", user.Id, user.Email)
		}
	})

	t.Run("Multiple entries", func(t *testing.T) {
		uTwo := models.User{Email: "email2@fake.com"}
		idTwo := createTestUser(&uTwo, t)
		uTwo.Id = idTwo
		usTwo := createTestSession(idTwo, t)

		user, err := db.GetUserBySessionToken(usTwo.User.Id, usTwo.SessionToken)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(user, uTwo) {
			t.Errorf("Expected %v to equal %v", user, uTwo)
		}
	})
}
