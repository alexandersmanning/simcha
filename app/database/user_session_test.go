package database

import (
	"testing"
	"github.com/alexandersmanning/simcha/app/models"
	"reflect"
)

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
