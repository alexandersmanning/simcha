package sessions

import (
	"errors"
	"net/http"

	"github.com/alexandersmanning/simcha/app/models"
	"github.com/gorilla/sessions"
)

type SessionStore interface {
	CurrentUser(r *http.Request) (*models.User, error)
	Login(u *models.User, w http.ResponseWriter, r *http.Request) error
	IsLoggedIn(r *http.Request) (bool, error)
	Logout(w http.ResponseWriter, r *http.Request) error
}

type Session struct {
	*sessions.CookieStore
}

func InitStore(secret string) *Session {
	cookieStore := sessions.NewCookieStore([]byte(secret))
	return &Session{ cookieStore }
}

func (s *Session) Login(u *models.User, w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["user"] = models.User{ID: u.ID, Email: u.Email}
	session.Values["loggedIn"] = true
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *Session) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["loggedIn"] = false
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *Session) IsLoggedIn(r *http.Request) (bool, error) {
	session, err := s.Get(r, "session")

	if err != nil {
		return false, err
	}

	if val, ok := session.Values["loggedIn"]; ok && val == true {
		return true, nil
	}

	return false, nil
}

func (s *Session) CurrentUser(r *http.Request) (*models.User, error) {
	// This does not work given scope values

	session, err := s.Get(r, "session")
	u := &models.User{}

	if err != nil {
		return u, err
	}

	if _, ok := session.Values["loggedIn"]; ok {
		val := session.Values["user"]
		if u, ok := val.(*models.User); !ok {
			return u, errors.New("session not stored properly")
		} else {
			return u, nil
		}
	}

	return u, nil
}
