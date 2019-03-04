package routes

import (
	"github.com/julienschmidt/httprouter"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/controllers"
	"github.com/alexandersmanning/simcha/app/middleware"
)

func Router(env *config.Env) *httprouter.Router {
	r := httprouter.New()
	r.GET("/posts", controllers.PostIndex(env))
	r.GET("/posts/:postId", controllers.PostIndex(env))
	r.POST("/posts", middleware.LoggedIn(env, controllers.PostCreate(env)))
	r.PUT("/posts/:postId", middleware.LoggedIn(env,
		middleware.PostPermission(env, controllers.PostUpdate(env)),
	))

	r.GET("/currentUser", controllers.CurrentUser(env))
	r.POST("/users", controllers.UserCreate(env))

	r.POST("/login", controllers.Login(env))
	r.GET("/logout", controllers.Logout(env))
	return r
}
