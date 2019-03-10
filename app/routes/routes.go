package routes

import (
	"github.com/julienschmidt/httprouter"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/controllers"
	"github.com/alexandersmanning/simcha/app/middleware"
)

func Router(env *config.Env) *httprouter.Router {
	r := httprouter.New()
	r.GET("/posts", middleware.CSRFMiddleware(
		env, controllers.PostIndex(env)),
	)
	r.GET("/posts/:postId", middleware.CSRFMiddleware(
		env, controllers.PostIndex(env)),
	)
	r.POST("/posts", middleware.LoggedIn(
		env, middleware.CSRFMiddleware(env, controllers.PostCreate(env))),
	)
	r.PUT("/posts/:postId", middleware.LoggedIn(env,
		middleware.PostPermission(
			env, middleware.CSRFMiddleware(
				env, controllers.PostUpdate(env),
			),
		),
	))
	r.DELETE("/posts/:postId", middleware.LoggedIn(env,
		middleware.PostPermission(
			env, middleware.CSRFMiddleware(
				env, controllers.PostDelete(env),
			),
		),
	))

	r.GET("/currentUser", middleware.CSRFMiddleware(
		env, controllers.CurrentUser(env),
	))
	r.POST("/users", middleware.CSRFMiddleware(
		env, controllers.UserCreate(env),
	))
	r.POST("/login", middleware.CSRFMiddleware(
		env, controllers.Login(env),
	))
	r.GET("/logout", middleware.CSRFMiddleware(
		env, controllers.Logout(env),
	))
	return r
}
