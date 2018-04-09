package sessions

import (
	"errors"
	"net/http"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
)

func Login(u *models.User, env *config.Env, w http.ResponseWriter, r *http.Request) error {
	session, err := env.Store.Get(r, "session")
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

func IsLoggedIn(env *config.Env, r *http.Request) (bool, error) {
	session, err := env.Store.Get(r, "session")

	if err != nil {
		return false, err
	}

	if val, ok := session.Values["loggedIn"]; ok && val == true {
		return true, nil
	}

	return false, nil
}

func CurrentUser(env *config.Env, r *http.Request) (*models.User, error) {
	session, err := env.Store.Get(r, "session")
	u := &models.User{}

	if err != nil {
		return u, err
	}

	if _, ok := session.Values["loggedIn"]; ok {
		val := session.Values["user"]
		if u, ok := val.(*models.User); !ok {
			return u, errors.New("Session not stored properly")
		}

		return u, nil
	}

	return u, nil
}
