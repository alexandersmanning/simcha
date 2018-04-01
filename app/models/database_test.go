package models

import (
	"database/sql"
	"fmt"
	"github.com/alexandersmanning/simcha/app/shared/database"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setupModels()
	exitCode := m.Run()

	teardownModels()
	os.Exit(exitCode)
}

func setupModels() {
	connString := "dbname=simcha_test sslmode=disable"
	db, err := sql.Open("postgres", connString)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	database.InitStore(db)
}

func teardownModels() {
	database.GetStore().Close()
}
