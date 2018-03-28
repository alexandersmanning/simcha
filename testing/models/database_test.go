package models

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/alexandersmanning/simcha/app/shared/database"
	_ "github.com/lib/pq"
)

type StoreSuite struct {
	suite.Suite
	store *database.DBStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	connString := "dbname=simcha_test sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}

	s.db = db
	database.InitStore(db)
	s.db = database.GetStore()
}

func (s *StoreSuite) TearDownSuite() {
	s.db.Close()
}

func (s *StoreSuite) SetupTest() {
	_, err := s.db.Query("DELETE FROM posts")
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.db.Query("DELETE FROM users")
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}
