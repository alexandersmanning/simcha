package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
)

func UserCreate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		var u models.User
		if err := json.Unmarshal(msg, &u); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		err = env.DB.CreateUser(&u)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		if err := env.Store.Login(&u, env.DB, w, r); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "applicaiton/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"result": "%v"}`, u.Email)
	}
}

func CurrentUser(env *config.Env) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u, err := env.Store.CurrentUser(env.DB, r)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(u)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}
