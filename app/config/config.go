package config

import (
	"github.com/alexandersmanning/simcha/app/sessions"
	"github.com/alexandersmanning/simcha/app/database"
)

type Env struct {
	DB    database.Datastore
	Store sessions.SessionStore
}
