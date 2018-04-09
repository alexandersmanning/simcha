package routes

import (
	"github.com/julienschmidt/httprouter"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/controllers"
)

func Router(env *config.Env) *httprouter.Router {
	r := httprouter.New()
	r.GET("/posts", controllers.PostIndex(env))
	r.POST("/posts", controllers.PostCreate(env))

	r.POST("/users", controllers.UserCreate(env))

	r.POST("/login", controllers.Login(env))
	r.GET("/logout", controllers.Logout(env))
	return r
}
