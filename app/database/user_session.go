package database

import (
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/alexandersmanning/webapputil"
)

type UserSessionStore interface {
	CreateUserSession(u *models.User) (models.UserSession, error)
	GetUserBySessionToken(userId int, token string) (models.User, error)
	RemoveSessionToken(userId int, token string) error
	RemoveAllUserSessions(userId int) error
}

func (db *DB) CreateUserSession(u *models.User) (models.UserSession, error) {
	var us models.UserSession

	token, err := CreateSessionToken()
	if err != nil {
		return us, err
	}

	rows, err := db.Query(`
		INSERT INTO user_sessions (user_id, session_token)
		VALUES ($1, $2) RETURNING id 
	`, u.Id, token)

	defer rows.Close()

	if err != nil {
		return us, err
	}

	for rows.Next() {
		if err := rows.Scan(&us.Id); err != nil {
			return us, err
		}
	}

	us.SessionToken, us.User = token, *u

	return us, nil
}

func (db *DB) GetUserBySessionToken(userId int, token string) (models.User, error) {
	var u models.User

	rows, err := db.Query(`
		SELECT DISTINCT users.id, users.email
		FROM users
		JOIN user_sessions ON (user_sessions.user_id = users.id)
		WHERE user_sessions.user_id = $1 AND user_sessions.session_token = $2
	`, userId, token)

	defer rows.Close()

	if err != nil {
		return u, err
	}

	for rows.Next() {
		if err := rows.Scan(&u.Id, &u.Email); err != nil {
			return u, err
		}
	}

	return u, nil
}

func (db *DB) RemoveSessionToken(userId int, token string) error {
	rows, err := db.Query(`
		DELETE FROM user_sessions WHERE user_id = $1 AND session_token = $2
	`, userId, token)

	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RemoveAllUserSessions(userId int) error {
	rows, err := db.Query(`
		DELETE FROM user_sessions WHERE user_id = $1
	`, userId)

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
