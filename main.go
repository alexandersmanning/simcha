package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"os"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/database"
	"github.com/alexandersmanning/simcha/app/routes"
	"github.com/alexandersmanning/simcha/app/sessions"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	db, err := database.InitDB(os.Getenv("DB_CONNECTION"))

	defer db.Close()

	if err != nil {
		panic(err)
	}

	store := sessions.InitStore(os.Getenv("APPLICATION_SECRET"))

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
	if err := http.ListenAndServe(":"+port, CorsHandler(csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(r))); err != nil {
		panic(err)
	}
}

func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin",  os.Getenv("DOMAIN"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Content-Length, X-CSRF-Token")
		h.ServeHTTP(w, r)
	})
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, r.URL.Path[1:]+"public")
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
