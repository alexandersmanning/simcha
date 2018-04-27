/*
Package models contains all of the applications models
and the Datastore and Generic functions and types to be used by all models
*/
package models

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	_ "github.com/lib/pq" //PQ is used for postgres db
)

//Datastore is the interface used by the router and mocks to interact with the database
type Datastore interface {
	PostStore
	UserStore
}

//DB is the public struct whose methods interact directly with the database
type DB struct {
	*sql.DB
}

//type M map[string]interface{}

//InitDB initializes the database, creating a new DB struct
func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

type modelError struct {
	fieldName string
	errorText string
}

func (m *modelError) Error() string {
	return fmt.Sprintf("%s %s", m.fieldName, m.errorText)
}

func init() {
	gob.Register(&User{})
}
