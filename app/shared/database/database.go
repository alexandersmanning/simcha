package database

import (
	"database/sql"
)

type DBStore struct {
	db *sql.DB
}

var store DBStore

func InitStore(s *sql.DB) {
	store.db = s
}

func GetStore() *sql.DB {
	return store.db
}
