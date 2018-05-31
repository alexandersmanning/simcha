package database

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/alexandersmanning/simcha/app/mocks/model"
)

func clearUsers(t *testing.T) {
	_, err := db.Query("DELETE FROM users")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUserExists(t *testing.T) {
	existsTest := func(input string, expected bool, t *testing.T) {
		t.Helper()
		exists, err := db.UserExists(input)

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

func TestCreateUser(t *testing.T) {
	email := "fake@email.com"
	password := "goodpassword"

	errorTestHelper := func(expectedName string, user models.User, t *testing.T) {
		t.Helper()

		if err := db.CreateUser(&user); err == nil {
			t.Error("Expected error, got nothing")
		} else if ae, ok := err.(*models.ModelError); !ok || ae.FieldName != expectedName {
			t.Errorf("Expected error with field %s, of %s", expectedName, err.Error())
		}
	}

	t.Run("User exists", func(t *testing.T) {
		clearUsers(t)

		if _, err := db.Query("INSERT INTO users (email) VALUES ($1)", email); err != nil {
			t.Fatal(err)
		}

		u := models.User{Email: email, Password: password, ConfirmationPassword: password}
		errorTestHelper("Email", u, t)
	})

	t.Run("Password is not the proper length", func(t *testing.T) {
		clearUsers(t)
		badPassword := "short"
		u := models.User{Email: email, Password: badPassword, ConfirmationPassword: badPassword}

		errorTestHelper("Password", u, t)
	})

	t.Run("Password does not match confirmation", func(t *testing.T) {
		clearUsers(t)
		badConfirmation := "nonmatchingpassword"
		u := models.User{Email: email, Password: password, ConfirmationPassword: badConfirmation}

		errorTestHelper("ConfirmationPassword", u, t)
	})

	t.Run("Creates a record and returns User with id", func(t *testing.T) {
		clearUsers(t)
		u := models.User{Email: email, Password: password, ConfirmationPassword: password}

		if err := db.CreateUser(&u); err != nil {
			t.Fatal(err)
		}

		if u.Id == 0 {
			t.Errorf("Id was not returned, expected int, got %d", u.Id)
		}
	})
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	clearUsers(t)

	email := "email@fake.com"
	password := "goodpassword"
	u := models.User{Email: email, Password: password, ConfirmationPassword: password}

	if err := db.CreateUser(&u); err != nil {
		t.Fatal(err)
	}

	t.Run("User exists in system", func(t *testing.T) {
		user, err := db.GetUserByEmailAndPassword(email, password)

		if err != nil {
			t.Error(err)
		}

		if user.Id != u.Id {
			t.Errorf("Expected to find user of ID %d, got %d", u.Id, user.Id)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		user, err := db.GetUserByEmailAndPassword("nonexistent@fake.com", password)

		if err == nil {
			t.Error("Expected error for missing user, got nothing")
		}

		if user.Id != 0 {
			t.Errorf("Expected no user to be found, got %v", user)
		}

		if ae, ok := err.(*models.ModelError); !ok || ae.FieldName != "Email or Password" {
			t.Errorf("Expected error %s, got %s", "Email or Password", err.Error())
		}
	})

	t.Run("User exists, password is incorrect", func(t *testing.T) {
		user, err := db.GetUserByEmailAndPassword(email, "wrongpassword")

		if err == nil {
			t.Error("Expected an erorr for wrong password, got nothing")
		}

		if user.Id != 0 {
			t.Errorf("Expected no user to be found, got %v", user)
		}

		if ae, ok := err.(*models.ModelError); !ok || ae.FieldName != "Email or Password" {
			t.Errorf("Expected %v error, got %v", "Email or Password", err)
		}
	})
}

func TestUpdatePassword(t *testing.T) {
	clearUsers(t)
	u := models.User{Email: "fake@email.com", Password: "correctPassword", ConfirmationPassword: "correctPassword" }
	if err := db.CreateUser(&u); err != nil {
		t.Fatal(err)
	}

	testHelper := func(t *testing.T,) {
		t.Helper()

	}

	t.Run("It fails if the previous password does not match", func(t *testing.T) {
		previousPassword, password, confirmation := "wrongPassword", "newPassword", "newPassword"
		expectedErr := &models.ModelError{FieldName: "Previous Password", ErrorText: "test text"}

		mockCtrl := gomock.NewController(t)
		userActions := mockmodel.NewMockUserAction(mockCtrl)

		userActions.EXPECT().ComparePassword(previousPassword).Return(expectedErr)

		err := db.UpdatePassword(
			 userActions,
			 previousPassword,
			 password,
			 confirmation,
		)

		if err == nil {
			t.Error("Expected an error for wrong password, got nothing")
		} else if val, ok := err.(*models.ModelError); !ok || val.FieldName != "Previous Password" {
			t.Errorf("Expected %v, got %v", "Previous Password", err)
		}
	})

	t.Run("It verifies the new passwords", func(t *testing.T) {
		previousPassword, password, confirmation := "rightPassword", "nonmatchpassword", "newPassword"
		expectedErr := &models.ModelError{FieldName: "Password", ErrorText: "Non Matching"}

		mockCtrl := gomock.NewController(t)
		mockUserAction := mockmodel.NewMockUserAction(mockCtrl)

		mockUserAction.EXPECT().ComparePassword(previousPassword).Return(nil).Times(1)
		mockUserAction.EXPECT().SetPassword(password, confirmation).Times(1)
		mockUserAction.EXPECT().CreateDigest().Return("", expectedErr).Times(1)

		err := db.UpdatePassword(mockUserAction, previousPassword, password, confirmation)

		if err == nil {
			t.Errorf("Expected an error for bad password match, got nothing")
		} else if val, ok := err.(*models.ModelError); !ok || val.FieldName != "Password" {
			t.Errorf("Expected %v, got %v", "Password", err)
		}
	})

}
