package database

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setupModels()
	exitCode := m.Run()

	teardownModels()
	os.Exit(exitCode)
}

var db *DB

func setupModels() {
	connString := "dbname=simcha_test sslmode=disable"
	var err error
	db, err = InitDB(connString)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func teardownModels() {
	db.Close()
}
