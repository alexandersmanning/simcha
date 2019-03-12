package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
)

func Login(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
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

		if err := env.Store.Login(&user, env.DB, w, r); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		jsonUser, err := json.Marshal(&user)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
		}

		sendJsonResponse(w, r, jsonUser)
	}
}

func Logout(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := env.Store.Logout(env.DB, w, r); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, r, "success")
	}
}
