package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/mocks/database"
	"github.com/alexandersmanning/simcha/app/mocks/sessions"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserCreate(t *testing.T) {
	u := models.User{Email: "email@fake.com", Password: "fakepassword", ConfirmationPassword: "fakepassword"}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mockdatabase.NewMockDatastore(mockCtrl)
	mockSessionStore := mocksession.NewMockSessionStore(mockCtrl)

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

		var resMsg JSONResponse
		err = json.Unmarshal(msg, &resMsg)

		if err != nil {
			t.Fatal(err)
		}

		if resMsg.Error != "failure" {
			t.Errorf("Expected an error, received %s", resMsg.Error)
		}
	})
}

func TestCurrentUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mockdatabase.NewMockDatastore(mockCtrl)
	mockSession := mocksession.NewMockSessionStore(mockCtrl)

	env := &config.Env{DB: mockDB, Store: mockSession}

	req, _ := http.NewRequest("GET", "/currentUser", nil)

	t.Run("Where this is a current user", func(t *testing.T) {
		res := httptest.NewRecorder()
		u := models.User{ Id: 123, Email: "email@fake.com" }
		mockSession.EXPECT().CurrentUser(mockDB, req).Return(&u, nil)
		CurrentUser(env)(res, req, nil)

		checkStatus(res.Code, 200, t)

		if res.Header().Get("Content-type") != "application/json" {
			t.Errorf("Expected to have json header, instead got %s", res.Header().Get("Content-type"))
		}

		var resUser models.User
		msg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(msg, &resUser)

		if err != nil {
			t.Fatal(err)
		}

		if resUser.Id != u.Id || resUser.Email != u.Email {
			t.Errorf("Expect %v go %v", u, resUser)
		}
	})

	t.Run("When there is no current user", func(t *testing.T) {
		res := httptest.NewRecorder()
		mockSession.EXPECT().CurrentUser(mockDB, req).Return(&models.User{}, nil)

		CurrentUser(env)(res, req, nil)

		checkStatus(res.Code, http.StatusOK, t)
		if res.Header().Get("Content-type") != "application/json" {
			t.Errorf("Expected application/json got %s", res.Header().Get("Content-type"))
		}

		var emptyUser models.User
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(resBody, &emptyUser)
		if err != nil {
			t.Fatal(err)
		}

		if emptyUser.Id != 0 || emptyUser.Email != "" {
			t.Errorf("Expected empty user, got %v", emptyUser)
		}
	})
}
