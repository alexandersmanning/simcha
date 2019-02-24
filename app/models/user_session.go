package models

type UserSession struct {
	Id int `json:"id"`
	User User `json:"user"`
	SessionToken string `json:"session_token"`
}
