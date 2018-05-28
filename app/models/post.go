package models

//Post is a struct for creating a simple blog post
type Post struct {
	Author User   `json:"author"`
	Body   string `json:"body"`
	Title  string `json:"title"`
}

