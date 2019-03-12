package controllers

import (
	"encoding/json"
	"errors"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func PostIndex(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		posts, err := env.DB.AllPosts()

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(posts)

		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		sendJsonResponse(w, r, body)
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

		sendJsonResponse(w, r, jsonPost)
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

		res := JSONResponse{ Result: "success" }
		jsonRes, err := json.Marshal(&res)
		sendJsonResponse(w, r, jsonRes)
	}
}

func PostDelete(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("postId")
		if id == "" {
			jsonError(w, errors.New("post Id must be provided"), http.StatusBadRequest)
			return
		}

		if err := env.DB.DeletePost(id); err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, r, "Success")
	}
}
