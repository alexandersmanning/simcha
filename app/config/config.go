package config

import (
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/gorilla/sessions"
)

type Env struct {
	DB    models.Datastore
	Store *sessions.CookieStore
}