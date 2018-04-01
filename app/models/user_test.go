package models

import (
	"testing"

	"github.com/alexandersmanning/simcha/app/shared/database"
)

func clearUsers(t *testing.T) {
	db := database.GetStore()
	_, err := db.Query("DELETE FROM users")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUserExists(t *testing.T) {
	db := database.GetStore()

	existsTest := func(input string, expected bool, t *testing.T) {
		t.Helper()
		exists, err := UserExists(input)

		if err != nil {
			t.Fatal(err)
		}

		if exists != expected {
			t.Errorf("Expected %v, got %v for %s", expected, exists, input)
		}
	}

	t.Run("No user available", func(t *testing.T) {
		clearUsers(t)
		existsTest("fake@email.com", false, t)
	})

	t.Run("Multiple users available", func(t *testing.T) {
		clearUsers(t)
		tests := []struct {
			input  string
			output bool
		}{
			{"email1@fake.com", true},
			{"email2@fake.com", true},
			{"email3@fake.com", false},
		}

		_, err := db.Query("INSERT INTO users (email) VALUES ($1), ($2)",
			"email1@fake.com", "email2@fake.com")

		if err != nil {
			t.Fatal(err)
		}

		for _, test := range tests {
			existsTest(test.input, test.output, t)
		}

	})
}

func TestEnsureSessionToken(t *testing.T) {
	clearUsers(t)
	u := User{Email: "email@fake.com"}

	if err := u.ensureSessionToken(); err != nil {
		t.Fatal(err)
	}

	token := u.SessionToken
	if token == "" {
		t.Errorf("Expected session token to be assigned, Session Token is %s", token)
	}

	if err := u.ensureSessionToken(); err != nil {
		t.Fatal(err)
	} else if u.SessionToken != token {
		t.Errorf("Expected session token to remain the %s, instead got %s", token, u.SessionToken)
	}
}
func TestCreateUser(t *testing.T) {
	db := database.GetStore()

	email := "fake@email.com"
	password := "goodpassword"

	errorTestHelper := func(expectedName string, user User, t *testing.T) {
		t.Helper()

		if err := user.CreateUser(); err == nil {
			t.Error("Expected error, got nothing")
		} else if ae, ok := err.(*modelError); !ok || ae.fieldName != expectedName {
			t.Errorf("Expected error with field %s, of %s", expectedName, err.Error())
		}
	}

	t.Run("User exists", func(t *testing.T) {
		clearUsers(t)

		if _, err := db.Query("INSERT INTO users (email) VALUES ($1)", email); err != nil {
			t.Fatal(err)
		}

		u := User{Email: email, Password: password, ConfirmationPassword: password}
		errorTestHelper("Email", u, t)
	})

	t.Run("Password is not the proper length", func(t *testing.T) {
		clearUsers(t)
		badPassword := "short"
		u := User{Email: email, Password: badPassword, ConfirmationPassword: badPassword}

		errorTestHelper("Password", u, t)
	})

	t.Run("Password does not match confirmation", func(t *testing.T) {
		clearUsers(t)
		badConfirmation := "nonmatchingpassword"
		u := User{Email: email, Password: password, ConfirmationPassword: badConfirmation}

		errorTestHelper("ConfirmationPassword", u, t)
	})

	t.Run("Creates a record and returns User with id", func(t *testing.T) {
		u := User{Email: email, Password: password, ConfirmationPassword: password}

		if err := u.CreateUser(); err != nil {
			t.Fatal(err)
		}

		if u.ID == 0 {
			t.Errorf("Id was not returned, expected int, got %d", u.ID)
		}

		if u.PasswordDigest == "" {
			t.Errorf("Expected a digest, got %s", u.PasswordDigest)
		}

		if u.SessionToken == "" {
			t.Errorf("Expected a session token, got %s", u.SessionToken)
		}
	})
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	clearUsers(t)

	email := "email@fake.com"
	password := "goodpassword"
	u := User{Email: email, Password: password, ConfirmationPassword: password}

	if err := u.CreateUser(); err != nil {
		t.Fatal(err)
	}

	t.Run("User exists in system", func(t *testing.T) {
		user, err := GetUserByEmailAndPassword(email, password)

		if err != nil {
			t.Error(err)
		}

		if user.ID != u.ID {
			t.Errorf("Expected to find user of ID %d, got %d", u.ID, user.ID)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		user, err := GetUserByEmailAndPassword("nonexistent@fake.com", password)

		if err == nil {
			t.Error("Expected error for missing user, got nothing")
		}

		if user.ID != 0 {
			t.Errorf("Expected no user to be found, got %v", user)
		}

		if ae, ok := err.(*modelError); !ok || ae.fieldName != "Email or Password" {
			t.Errorf("Expected error %s, got %s", "Email or Password", err.Error())
		}
	})

	t.Run("User exists, password is incorrect", func(t *testing.T) {
		user, err := GetUserByEmailAndPassword(email, "wrongpassword")

		if err == nil {
			t.Error("Expected an erorr for wrong password, got nothing")
		}

		if user.ID != 0 {
			t.Errorf("Expected no user to be found, got %v", user)
		}

		if ae, ok := err.(*modelError); !ok || ae.fieldName != "Email or Password" {
			t.Errorf("Expected %v error, got %v", "Email or Password", err)
		}
	})
}
