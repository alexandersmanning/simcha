package routes

import (
	"github.com/julienschmidt/httprouter"

	"github.com/alexandersmanning/simcha/app/controllers"
)

func Router() *httprouter.Router {
	r := httprouter.New()
	r.GET("/posts", controllers.PostIndex)
	r.POST("/posts", controllers.PostCreate)
	return r
}
