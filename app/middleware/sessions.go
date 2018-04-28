package middleware

import (
	"github.com/julienschmidt/httprouter"
	"github.com/alexandersmanning/simcha/app/config"
	"net/http"
)

type Middleware func(next httprouter.Handle) httprouter.Handle

func LoggedIn(env *config.Env, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		if loggedIn, err := env.Store.IsLoggedIn(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if !loggedIn {
			http.Error(w, "You must be logged in", http.StatusInternalServerError)
			return
		}

		next(w, r, param)
	}
}
