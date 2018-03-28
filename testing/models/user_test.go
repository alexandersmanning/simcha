package models

import (
	"github.com/alexandersmanning/simcha/app/models"
)

func (s *StoreSuite) TestUserExists() {
	exists, err := models.UserExists("email@fake.com")
	if err != nil {
		s.T().Fatal(err)
	} else if exists {
		s.T().Errorf("Expected no users to exists, instead received error %v\n", err)
	}

	//create user and then check
	_, err = s.db.Query(
		"INSERT INTO users (email) VALUES ($1), ($2)",
		"email1@fake.com", "email2@fake.com")

	if err != nil {
		s.T().Fatal(err)
	}

	testing := []struct {
		input  string
		output bool
	}{
		{"email1@fake.com", true},
		{"email2@fake.com", true},
		{"email3@fake.com", false},
	}

	for _, test := range testing {
		exists, err = models.UserExists(test.input)
		if err != nil {
			s.T().Fatal(err)
		} else if exists != test.output {
			s.T().Errorf("Expected %v to be %v, instead got %v", test.input, test.output, exists)
		}
	}
}

func (s *StoreSuite) TestCreatePassword() {
}
