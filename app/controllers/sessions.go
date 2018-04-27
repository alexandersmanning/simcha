package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/alexandersmanning/simcha/app/sessions"
)

func Login(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			jsonError(w,err, http.StatusInternalServerError)
			return
		}

		var user models.User
		if err := json.Unmarshal(msg, &user); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		user, err = env.DB.GetUserByEmailAndPassword(user.Email, user.Password)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		if err := sessions.Login(&user, env, w, r); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"result": "success"}`)
	}
}

func Logout(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		session, err := env.Store.Get(r, "session")
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		session.Values["loggedIn"] = false
		session.Save(r, w)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"result": "success"}`)
	}
}
