package sessions

import (
	"net/http"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/gorilla/sessions"
	"errors"
)

type SessionStore interface {
	CurrentUser(db models.Datastore, r *http.Request) (*models.User, error)
	Login(u *models.User, w http.ResponseWriter, r *http.Request) error
	IsLoggedIn(db models.Datastore, r *http.Request) (bool, error)
	Logout(u *models.User, db models.Datastore, w http.ResponseWriter, r *http.Request) error
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

	session.Values["id"] = u.ID
	session.Values["token"] = u.SessionToken

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *Session) Logout(u *models.User, db models.Datastore, w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["token"] = nil
	session.Values["id"] = nil
	if err := session.Save(r, w); err != nil {
		return err
	}

	// User does not exists, just return nil
	if u.ID == 0 {
		return nil
	}

	return db.UpdateSessionToken(u.ID)
}

func (s *Session) IsLoggedIn(db models.Datastore, r *http.Request) (bool, error) {
	u, err := s.CurrentUser(db, r)

	if err != nil {
		return false, err
	}

	if u.ID == 0 {
		return false, nil
	}

	return true, nil
}

func (s *Session) CurrentUser(db models.Datastore, r *http.Request) (*models.User, error) {
	var u models.User

	id, token, err := getSessionValues(s, r)
	if err != nil {
		return u, err
	}

	u, err = db.GetUserBySessionToken(id, token)
	return &u, err
}

func getSessionValues(s *Session, r *http.Request) (int, string, error){
	session, err := s.Get(r, "session")

	var id int
	var token string

	if err != nil {
		return id, token, err
	}

	val := session.Values["id"]
	id, ok := val.(int)

	if !ok {
		return id, token, errors.New("no session id found")
	}

	val = session.Values["token"]
	token, ok = val.(string)

	if !ok {
		return id, token, errors.New("no session token found")
	}


	return  id, token, nil

}
