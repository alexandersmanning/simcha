package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
)

func PostIndex(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		posts, err := env.DB.AllPosts()

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		err = enc.Encode(posts)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func PostCreate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		var post models.Post
		err = json.Unmarshal(msg, &post)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		user, err := env.Store.CurrentUser(env.DB, r)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		post.Author = *user

		err = env.DB.CreatePost(&post)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}


		jsonPost, err := json.Marshal(&post)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonPost)
	}
}

func PostUpdate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		var post models.Post
		err = json.Unmarshal(bytes, &post)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		err = env.DB.EditPost(&post)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, "success")
	}
}
