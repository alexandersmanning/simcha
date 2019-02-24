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
