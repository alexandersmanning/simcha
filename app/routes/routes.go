package routes

import (
	"github.com/julienschmidt/httprouter"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/controllers"
	"github.com/alexandersmanning/simcha/app/middleware"
)

func Router(env *config.Env) *httprouter.Router {
	r := httprouter.New()
	r.GET("/posts", middleware.LoggedIn(env, controllers.PostIndex(env)))
	r.POST("/posts", middleware.LoggedIn(env, controllers.PostCreate(env)))

	r.POST("/users", controllers.UserCreate(env))

	r.POST("/login", controllers.Login(env))
	r.GET("/logout", controllers.Logout(env))
	return r
}
