package models

import "time"

//Post is a struct for creating a simple blog post
type Post struct {
	Id         int       `json:"id"`
	Author     User      `json:"author"`
	Body       string    `json:"body"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	ModifiedAt time.Time `json:"updatedAt,omitempty"`
}

type PostAction interface {
	ModelAction
	Post() *Post
}

func (p *Post) Post() *Post {
	return p
}

func (p *Post) SetID(id int) {
	p.Id = id
}

func (p *Post) SetTimestamps() {
	if (p.CreatedAt == time.Time{}) {
		p.CreatedAt = time.Now().UTC()
	}

	p.ModifiedAt = time.Now().UTC()
}

func (p *Post) Timestamps() (time.Time, time.Time) {
	return p.CreatedAt, p.ModifiedAt
}
