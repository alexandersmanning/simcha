package sessions

import (
	"net/http"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/gorilla/sessions"
	"errors"
	"github.com/alexandersmanning/simcha/app/database"
)

type SessionStore interface {
	CurrentUser(db database.Datastore, r *http.Request) (*models.User, error)
	Login(u *models.User, db database.Datastore, w http.ResponseWriter, r *http.Request) error
	IsLoggedIn(db database.Datastore, r *http.Request) (bool, error)
	Logout(db database.Datastore, w http.ResponseWriter, r *http.Request) error
}

type Session struct {
	*sessions.CookieStore
}

func InitStore(secret string) *Session {
	cookieStore := sessions.NewCookieStore([]byte(secret))
	return &Session{ cookieStore }
}


func (s *Session) Login(u *models.User, db database.Datastore, w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	us, err := db.CreateUserSession(u)
	if err != nil {
		return err
	}

	session.Values["id"] = us.Id
	session.Values["token"] = us.SessionToken

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *Session) Logout(db database.Datastore, w http.ResponseWriter, r *http.Request) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	currentId, currentToken, err := getSessionValues(s, r)
	if err != nil {
		return err
	}

	session.Values["token"] = nil
	session.Values["id"] = nil
	if err := session.Save(r, w); err != nil {
		return err
	}

	// User does not exists, just return nil
	if currentId == 0 {
		return nil
	}

	return db.RemoveSessionToken(currentId, currentToken)
}

func (s *Session) IsLoggedIn(db database.Datastore, r *http.Request) (bool, error) {
	u, err := s.CurrentUser(db, r)

	if err != nil {
		return false, err
	}

	if u.Id == 0 {
		return false, nil
	}

	return true, nil
}

func (s *Session) CurrentUser(db database.Datastore, r *http.Request) (*models.User, error) {
	var u models.User

	id, token, err := getSessionValues(s, r)
	if err != nil {
		return &u, err
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
