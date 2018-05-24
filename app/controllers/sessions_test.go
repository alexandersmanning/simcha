package controllers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/golang/mock/gomock"
	"github.com/alexandersmanning/simcha/app/mocks"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
)

func TestLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSessionStore := mocks.NewMockSessionStore(mockCtrl)
	mockDataStore := mocks.NewMockDatastore(mockCtrl)
	env := config.Env{DB: mockDataStore, Store: mockSessionStore}
	u := models.User{Email: "fake@email.com", Password: "thisisatestpassword"}

	t.Run("Matching credentials", func(t *testing.T) {
		jsonUser, err := json.Marshal(u)
		if err != nil {
			t.Fatal(err)
		}
		userBuff := bytes.NewBuffer(jsonUser)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", userBuff)
		mockDataStore.EXPECT().GetUserByEmailAndPassword(u.Email, u.Password).Return(u, nil)
		mockSessionStore.EXPECT().Login(&u, env.DB, rec, req).Return(nil)

		Login(&env)(rec, req, nil)

		checkStatus(rec.Code, 200, t)
		checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

		msg, err := ioutil.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(t)
		}

		if string(msg) != `{"result": "success"}` {
			t.Errorf("Expected successful result, got %s", string(msg))
		}
	})

	t.Run("No user found", func(t *testing.T) {
		jsonUser, err := json.Marshal(u)
		if err != nil {
			t.Fatal(err)
		}

		userBuff := bytes.NewBuffer(jsonUser)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", userBuff)

		mockDataStore.EXPECT().GetUserByEmailAndPassword(u.Email, u.Password).Return(models.User{}, errors.New("no user found"))

		Login(&env)(rec, req, nil)

		checkStatus(rec.Code, 500, t)
		checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

		msg, err := ioutil.ReadAll(rec.Body)

		if err != nil {
			t.Fatal(t)
		}

		if string(msg) != `{"result":"", "error": "no user found"}` {
			t.Errorf("Expected an error, received %s instead", string(msg))
		}
	})
}

func TestLogout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDataStore := mocks.NewMockDatastore(mockCtrl)
	mockSessionStore := mocks.NewMockSessionStore(mockCtrl)

	env := &config.Env{DB: mockDataStore, Store: mockSessionStore}

	t.Run("Working Session Logout", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/logout", nil)
		rec := httptest.NewRecorder()

		mockSessionStore.EXPECT().Logout(env.DB, rec, req).Return(nil)
		Logout(env)(rec, req, nil)

		checkStatus(rec.Code, 200, t)
		checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

		msg, err := ioutil.ReadAll(rec.Body)

		if err != nil {
			t.Fatal(err)
		}

		if string(msg) != `{"result": "success"}` {
			t.Errorf("Expected positive result, receive %s", string(msg))
		}
	})

	t.Run("Failed Session Logout", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/logout", nil)
		rec := httptest.NewRecorder()

		mockSessionStore.EXPECT().Logout(env.DB, rec, req).Return(errors.New("session failed"))

		Logout(env)(rec, req, nil)

		checkStatus(rec.Code, 500, t)
		checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

		msg, err := ioutil.ReadAll(rec.Body)

		if err != nil {
			t.Fatal(err)
		}

		if string(msg) != `{"result":"", "error": "session failed"}` {
			t.Errorf("expected error result, got %s", string(msg))
		}
	})
}
