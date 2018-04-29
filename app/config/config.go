package config

import (
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/gorilla/sessions"
	"net/http"
)

type Env struct {
	DB    models.Datastore
	Store *sessions.CookieStore
}

type EnvStore interface {
	CurrentUser(r *http.Request) (*models.User, error)
	Login(u *models.User, w http.ResponseWriter, r *http.Request) error
	IsLoggedIn(r *http.Request) (bool, error)
	Logout(w http.ResponseWriter, r *http.Request) error
}


func InitStore(secret string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(secret))
}

func (e *Env) Login(u *models.User, w http.ResponseWriter, r *http.Request) error {
	session, err := e.Store.Get(r, "session")
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

func (e *Env) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := e.Store.Get(r, "session")
	if err != nil {
		return err
	}

	//get current user and then update
	u, err := e.CurrentUser(r)
	if err != nil {
		return err
	}

	if u.ID == 0 {
		return nil
	}

	session.Values["token"] = nil
	session.Values["id"] = nil
	if err := session.Save(r, w); err != nil {
		return err
	}

	return e.DB.UpdateSessionToken(u.ID)
}

func (e *Env) IsLoggedIn(r *http.Request) (bool, error) {
	u, err := e.CurrentUser(r)

	if err != nil {
		return false, err
	}

	if u.ID == 0 {
		return false, nil
	}

	return true, nil
}

func (e *Env) CurrentUser(r *http.Request) (*models.User, error) {
	var u models.User

	session, err := e.Store.Get(r, "session")

	if err != nil {
		return &u, err
	}

	val := session.Values["id"]
	id, ok := val.(int);

	if !ok {
		return &u, nil
	}

	val = session.Values["token"]
	token, ok := val.(string)

	if !ok {
		return &u, nil
	}

	u, err = e.DB.GetUserBySessionToken(id, token)
	return &u, err
}
