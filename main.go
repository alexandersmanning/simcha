package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/alexandersmanning/simcha/app/routes"
	"github.com/alexandersmanning/simcha/app/sessions"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	db, err := models.InitDB(os.Getenv("DB_CONNECTION"))

	defer db.Close()

	if err != nil {
		panic(err)
	}

	store := sessions.InitStore(os.Getenv("APPLICATION_SECRET")) //.NewCookieStore([]byte(os.Getenv("APPLICATION_SECRET")))

	env := &config.Env{DB: db, Store: store}
	r := routes.Router(env)

	//// this is a generic serve for things like CSS
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	r.GET("/", Index)
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/404.html")
	})

	port := os.Getenv("PORT")
	fmt.Println("Listening on ", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, r.URL.Path[1:] + "public")
	//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//	dat, err := ioutil.ReadFile("public/index.html")
	//	if err == nil {
	//		fmt.Fprintf(w, string(dat))
	//	} else {
	//		fmt.Println(err)
	//	}
}

//func PostsIndexHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	post := Post{Author: "Alex Manning", Title: "This is my post", Body: "Posts for Days"}
//	enc := json.NewEncoder(w)
//	//js, err := json.Marshal(post)
//	//if err != nil {
//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
//	//	return
//	//}
//	w.Header().Set("Content-Type", "application/json")
//	err := enc.Encode(post)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	//fmt.Fprintln(w, string(js))
//}
