package config

import (
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/alexandersmanning/simcha/app/sessions"
)

type Env struct {
	DB    models.Datastore
	Store sessions.SessionStore
}
