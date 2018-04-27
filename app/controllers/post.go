package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
	//"github.com/alexandersmanning/simcha/app/sessions"
)

func PostIndex(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//if loggedIn, err := sessions.IsLoggedIn(env, r); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//} else if !loggedIn {
		//	http.Error(w, "You must be logged in", http.StatusInternalServerError)
		//	return
		//}

		w.Header().Set("Content-Type", "application/json")
		posts, err := env.DB.AllPosts()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)

		err = enc.Encode(posts)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PostCreate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var post models.Post
		err = json.Unmarshal(msg, &post)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.DB.CreatePost(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{"result": "success"}`)
	}
}
