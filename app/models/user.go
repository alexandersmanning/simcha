package models

import (
	"crypto/sha256"
	"errors"
	"github.com/alexandersmanning/simcha/app/shared/database"
	//"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"hash"
	"time"
)

// User struct is exported
type User struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
	PasswordDigest       string
	SessionToken         hash.Hash
	CreatedAt            time.Time `json:"createdAt"`
	ModifiedAt           time.Time `json:"modifiedAt"`
}

func (u *User) ensureSessionToken() {
	token := sha256.New()
	if u.SessionToken == nil {
		u.SessionToken = token
	}
}

func (u *User) createPassword() error {
	err := verifyPassword(u.Password, u.ConfirmationPassword)
	if err != nil {
		return err
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}

	u.PasswordDigest = string(passwordByte)
	return nil
}

func verifyPassword(password, confirmation string) error {
	if len(password) < 6 {
		return errors.New("Password must be at least 6 characters")
	}

	if password != confirmation {
		return errors.New("Password does not match confirmation")
	}

	return nil
}

func UserExists(email string) (bool, error) {
	var count int
	rows, err := database.GetStore().Query("SELECT COUNT(*) FROM users WHERE email = $1", email)

	if err != nil {
		return false, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}

		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (u *User) createUser() (User, error) {
	// verify email and users does not exists
	exists, err := UserExists(u.Email)

	if err != nil {
		return *u, err
	}

	if exists {
		return *u, errors.New("User already exists in the system")
	}
	// set password
	u.createPassword()
	u.ensureSessionToken()

	return *u, nil
}
