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

func UserCreate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var u models.User
		if err := json.Unmarshal(msg, &u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.DB.CreateUser(&u)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := sessions.Login(&u, env, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "applicaiton/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"result": "%v"}`, u.Email)
	}
}
