package models

import "github.com/alexandersmanning/webapputil"

type UserSession struct {
	ID int `json:"id"`
	User User `json:"user"`
	SessionToken string `json:"token"`
}

type UserSessionStore interface {
	CreateUserSession(u *User) (UserSession, error)
	GetUserBySessionToken(userId int, token string) (User, error)
}

func (db *DB) CreateUserSession(u *User) (UserSession, error){
	var us UserSession

	token, err := CreateSessionToken()
	if err != nil {
		return us, err
	}

	rows, err := db.Query(`
		INSERT INTO user_sessions (user_id, token)
		VALUES ($1, $2) RETURNING id 
	`, u.ID, token)

	defer rows.Close()

	if err != nil {
		return us, err
	}

	for rows.Next() {
		if err := rows.Scan(&us.ID); err != nil {
			return us, err
		}
	}

	return us, nil
}

func (db *DB) GetUserBySessionToken(userId int, token string) (User, error){
	var u User

	rows, err := db.Query(`
		SELECT DISTINCT users.id, users.email
		FROM users
		JOIN user_sessions ON (user_sessions.user_id = users.id)
		WHERE user_sessions.user_id = $1 AND user_sessions.token = $2
	`, userId, token)

	defer rows.Close()

	if err != nil {
		return u, err
	}

	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Email); err != nil {
			return u, err
		}
	}

	return u, nil
}

func (db *DB) RemoveSessionToken(userId int, token string) (error) {
	rows, err := db.Query(`
		DELETE FROM user_sessions WHERE user_id = $1 AND token = $2
	`, userId, token)

	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func CreateSessionToken() (string, error) {
	token, err := webapputil.GenerateSecureRandom()
	if err != nil {
		return "", err
	}

	return token, nil
}
