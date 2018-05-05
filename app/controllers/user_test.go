package controllers

import (
	"testing"
	"net/http"
	"github.com/alexandersmanning/simcha/app/models"
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/alexandersmanning/simcha/app/mocks"
	"github.com/alexandersmanning/simcha/app/config"
	"errors"
	"io/ioutil"
)

func TestUserCreate(t *testing.T) {
	u := models.User{Email: "email@fake.com", Password: "fakepassword", ConfirmationPassword: "fakepassword"}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	mockSessionStore := mocks.NewMockSessionStore(mockCtrl)

	env := &config.Env{DB: mockDatastore, Store: mockSessionStore}

	jsonUser, err := json.Marshal(u)

	if err != nil {
		t.Fatal(err)
	}
	userBuff := bytes.NewBuffer(jsonUser)

	req, _ := http.NewRequest("POST", "/users", userBuff)
	rec := httptest.NewRecorder()

	t.Run("Failure creating user", func(t *testing.T) {
		mockDatastore.EXPECT().CreateUser(&u).Return(errors.New("failure"))
		UserCreate(env)(rec, req, nil)

		checkStatus(rec.Code, 500, t)
		checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

		msg, err := ioutil.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(msg) != `{"result":"", "error": "failure"}` {
			t.Errorf("Expected an error, received %s", string(msg))
		}
	})
}
