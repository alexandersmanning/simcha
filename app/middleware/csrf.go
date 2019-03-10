package middleware

import (
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func CSRFMiddleware(_ *config.Env, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		next(w, r, p)
	}
}
