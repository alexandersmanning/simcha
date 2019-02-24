package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/mocks/database"
	"github.com/alexandersmanning/simcha/app/mocks/sessions"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostPermission(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockDB := mockdatabase.NewMockDatastore(mockCtrl)
	mockStore := mocksession.NewMockSessionStore(mockCtrl)

	env := config.Env{DB: mockDB, Store: mockStore}

	calledMockFunc := false
	mockFunc := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		calledMockFunc = true
	}

	user := models.User{Email: "email@fake.com", Id: 1 }
	otherUser := models.User{Email: "other@fake.com", Id: 3 }
	post := models.Post{Id: 2, Body: "fakeBody", Title: "fakeTitle", Author: user }
	updateBody, err := json.Marshal(models.Post{Id: 2, Body: "UpdatedBody"})
	buffBody := bytes.NewBuffer(updateBody)

	if err != nil {
		t.Fatal(err.Error())
	}

	req, _ := http.NewRequest("PUT", "/posts/2", buffBody)
	params := []httprouter.Param{{Key: "postId", Value: "2"}}

	t.Run("It calls an error if the DB cannot be called", func (t *testing.T) {
		res := httptest.NewRecorder()
		mockDB.EXPECT().GetPostById("2").Return(nil, errors.New("failure"))
		PostPermission(&env, mockFunc)(res, req, params)
		if res.Code != 500 {
			t.Errorf("Expected to receive 500, got %d", res.Code)
		}

		if calledMockFunc == true {
			t.Error("Expected next not to have been called")
		}
	})

	t.Run("Current User returns an error", func(t *testing.T) {
		res := httptest.NewRecorder()
		mockDB.EXPECT().GetPostById("2").Return(&post, nil)
		mockStore.EXPECT().CurrentUser(env.DB, req).Return(nil, errors.New("failure"))

		PostPermission(&env, mockFunc)(res, req, params)
		if res.Code != 500 {
			t.Errorf("Expected to receive a code of 500, got %d", res.Code)
		}

		if calledMockFunc == true {
			t.Error("Expected next not to have been called")
		}
	})

	t.Run("Current user does not match post user", func(t *testing.T) {
		res := httptest.NewRecorder()

		calledMockFunc = false
		mockStore.EXPECT().CurrentUser(env.DB, req).Return(&otherUser,nil)
		mockDB.EXPECT().GetPostById("2").Return(&post, nil)

		PostPermission(&env, mockFunc)(res, req, params)

		if res.Code != 400 {
			t.Errorf("Expected to get a 400 code, got %d", res.Code)
		}

		if calledMockFunc == true {
			t.Error("Expected next not to have been called")
		}
	})

	t.Run("Current user matches ", func(t *testing.T) {
		res := httptest.NewRecorder()

		calledMockFunc = false
		mockStore.EXPECT().CurrentUser(env.DB, req).Return(&user,nil)
		mockDB.EXPECT().GetPostById("2").Return(&post, nil)

		PostPermission(&env, mockFunc)(res, req, params)

		if calledMockFunc != true {
			t.Error("Expected next to have been called")
		}
	})
}
