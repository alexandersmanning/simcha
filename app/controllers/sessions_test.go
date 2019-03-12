package controllers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/golang/mock/gomock"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
	"github.com/alexandersmanning/simcha/app/mocks/database"
	"github.com/alexandersmanning/simcha/app/mocks/sessions"
)

func TestLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSessionStore := mocksession.NewMockSessionStore(mockCtrl)
	mockDataStore := mockdatabase.NewMockDatastore(mockCtrl)
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

		uFound := models.User{}
		err = json.Unmarshal(msg, &uFound)
		if err != nil {
			t.Fatal(err)
		}

		if u.Email != uFound.Email || u.Id != uFound.Id {
			t.Errorf("Expected %v, got %v", u, uFound)
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

		var bodyRes JSONResponse
		err = json.Unmarshal(msg, &bodyRes)
		if err != nil {
			t.Fatal(t)
		}

		if bodyRes.Error != "no user found" {
			t.Errorf("Expected an error, received %s instead", bodyRes.Error)
		}
	})
}

func TestLogout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDataStore := mockdatabase.NewMockDatastore(mockCtrl)
	mockSessionStore := mocksession.NewMockSessionStore(mockCtrl)

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

		jsonRes := JSONResponse{};
		if err := json.Unmarshal(msg, &jsonRes); err != nil {
			t.Fatal(err)
		}

		if jsonRes.Result != "success" {
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

		var resMsg JSONResponse
		err = json.Unmarshal(msg, &resMsg)
		if resMsg.Error != "session failed" {
			t.Errorf("expected error result, got %s", resMsg.Error)
		}
	})
}
