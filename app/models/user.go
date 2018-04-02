package models

import (
	"time"

	//"github.com/gorilla/sessions"
	"github.com/alexandersmanning/webapputil"
	"golang.org/x/crypto/bcrypt"
)

// User struct is exported
type User struct {
	ID                   int    `json:"id"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
	PasswordDigest       string
	SessionToken         string
	CreatedAt            time.Time `json:"createdAt"`
	ModifiedAt           time.Time `json:"modifiedAt"`
}

//UserStore is the interface for all User functions that interact with the database
type UserStore interface {
	GetUserByEmailAndPassword(email, password string) (User, error)
	UserExists(email string) (bool, error)
	CreateUser(u *User) error
}

//GetUserByEmailAndPassword checks if the user is in the database, and if it is verifies if the password matches
func (db *DB) GetUserByEmailAndPassword(email, password string) (User, error) {
	u := User{}
	rows, err := db.Query(
		`SELECT id, email, password_digest FROM users WHERE email = $1`, email,
	)

	if err != nil {
		return User{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordDigest); err != nil {
			return User{}, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))

	if u.Email == "" || err != nil {
		return User{}, &modelError{"Email or Password", "was not found, or does not match our records"}
	}

	return u, nil
}

func (u *User) createPassword() error {
	err := verifyPassword(u.Password, u.ConfirmationPassword)
	if err != nil {
		return err
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordDigest = string(passwordByte)
	return nil
}

func verifyPassword(password, confirmation string) error {
	if len(password) < 6 {
		return &modelError{"Password", "must be at least 6 characters long"}
	}

	if password != confirmation {
		return &modelError{"ConfirmationPassword", "does not match Password"}
	}

	return nil
}

//UserExists checks the existence of an email
func (db *DB) UserExists(email string) (bool, error) {
	var count int
	rows, err := db.Query("SELECT COUNT(*) FROM users WHERE email = $1", email)

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

//CreateUser adds user to system if they do not already exist, and have an appropriate email/password
func (db *DB) CreateUser(u *User) error {
	if exists, err := db.UserExists(u.Email); err != nil {
		return err
	} else if exists {
		return &modelError{"Email", "already exists in the system"}
	}

	// set password
	if err := u.createPassword(); err != nil {
		return err
	}

	if err := u.ensureSessionToken(); err != nil {
		return err
	}

	//to be handled by the middleware
	u.CreatedAt = time.Now().UTC()
	u.ModifiedAt = time.Now().UTC()

	rows, err := db.Query(`
		INSERT INTO users (email, password_digest, session_token, created_at, modified_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, u.Email, u.PasswordDigest, u.SessionToken, u.CreatedAt, u.ModifiedAt)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&u.ID); err != nil {
			return err
		}
	}

	return nil
}

func (u *User) ensureSessionToken() error {
	if token, err := CreateSessionToken(); err != nil {
		return err
	} else if u.SessionToken == "" {
		u.SessionToken = token
	}

	return nil
}

//CreateSessionToken returns an base64 secure random token
func CreateSessionToken() (string, error) {
	token, err := webapputil.GenerateSecureRandom()

	if err != nil {
		return "", err
	}

	return token, nil
}
