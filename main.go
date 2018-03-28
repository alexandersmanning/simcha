package main

import (
	//	"encoding/json"
	"fmt"
	// "github.com/julienschmidt/httprouter"
	//	"github.com/russross/blackfriday"
	//	"io/ioutil"
	"database/sql"
	"net/http"

	"github.com/alexandersmanning/simcha/app/routes"
	"github.com/alexandersmanning/simcha/app/shared/database"
	_ "github.com/lib/pq"
)

func main() {
	r := routes.Router()
	db, err := sql.Open("postgres", "dbname=simcha_dev sslmode=disable")
	database.InitStore(db)

	if err != nil {
		panic(err)
	}
	//r := httprouter.New()
	//r.GET("/", HomeHandler)

	//// post collection
	//r.GET("/posts", PostsIndexHandler)
	//r.POST("/posts", PostsCreateHandler)

	//// Posts singular
	//r.GET("/posts/:id", PostShowHandler)
	//r.PUT("/posts/:id", PostUpdateHandler)
	//r.GET("/posts/:id/edit", PostEditHandler)

	//// this is a generic serve for things like CSS
	//r.ServeFiles("/public/*filepath", http.Dir("public"))
	//r.POST("/markdown", GenerateMarkdown)
	////http.HandleFunc("/markdown", GenerateMarkdown)
	////http.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Println("Listening on ", 8080)
	if error := http.ListenAndServe(":8080", r); error != nil {
		panic(error)
	}
}

//func GenerateMarkdown(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	markdown := blackfriday.MarkdownCommon([]byte(r.FormValue("body")))
//	w.Write(markdown)
//}
//
//func HomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	// this purpoesley sets the headers, reads the files, and outputs them with the response
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	dat, err := ioutil.ReadFile("public/index.html")
//	if err == nil {
//		fmt.Fprintf(w, string(dat))
//	} else {
//		fmt.Println(err)
//	}
//}
//
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
//
//func PostsCreateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	fmt.Fprintln(w, "get create")
//}
//
//func PostShowHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	id := p.ByName("id")
//	fmt.Fprintln(w, "showing posts: ", id)
//}
//
//func PostUpdateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	id := p.ByName("id")//
//	fmt.Fprintln(w, "editing post: ", id)
//}

//func PostEditHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	fmt.Fprintln(w, "post again")
//}
