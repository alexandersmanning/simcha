package middleware

import (
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func PostPermission(env *config.Env, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postId := p.ByName("postId")
		post, err := env.DB.GetPostById(postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := env.Store.CurrentUser(env.DB, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user.Id != post.Author.Id || user.Email != post.Author.Email {
			http.Error(w, "User does not match", http.StatusBadRequest)
			return
		}
		next(w, r, p)
	}
}
