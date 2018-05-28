package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct is exported
type User struct {
	Id                   int    `json:"id"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
	PasswordDigest       string	`json:"-"`
	CreatedAt            time.Time `json:"createdAt,omitempty"`
	ModifiedAt           time.Time `json:"modifiedAt,omitempty"`
}

// Idea from https://stackoverflow.com/questions/26027350/go-interface-fields
type UserAction interface {
	ModelAction
	User() *User
	SetPassword(password, confirmation string)
	SetDigest(digest string)
	CreateDigest() (string, error)
	VerifyPassword() error
	ComparePassword(password string) error
}

func (u *User) User() *User {
	return u
}

func (u *User) CreateDigest() (string, error) {
	var digest string
	if err := u.VerifyPassword(); err != nil {
		return digest, err
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return digest, err
	}

	digest = string(passwordByte)
	return digest, nil
}

func (u *User) VerifyPassword() error {
	if len(u.Password) < 6 {
		return &ModelError{"Password", "must be at least 6 characters long"}
	}

	if u.Password != u.ConfirmationPassword {
		return &ModelError{"ConfirmationPassword", "does not match Password"}
	}

	return nil
}

func (u *User) Timestamps() (time.Time, time.Time) {
	return u.CreatedAt, u.ModifiedAt
}

func (u *User) SetTimestamps() {
	if (u.CreatedAt == time.Time{}) {
		u.CreatedAt = time.Now().UTC()
	}
	u.ModifiedAt = time.Now().UTC()
}

func (u *User) SetID(id int) {
	u.Id = id
}

func (u *User) SetPassword(password, confirmation string) {
	u.Password, u.ConfirmationPassword = password, confirmation
}

func (u *User) SetDigest(digest string) {
	u.PasswordDigest = digest
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password)); err != nil {
		return err
	}

	return nil
}

