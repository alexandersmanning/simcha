package controllers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/golang/mock/gomock"
	"github.com/alexandersmanning/simcha/app/mocks"
	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/models"
)

func TestLogin(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSessionStore := mocks.NewMockSessionStore(mockCtrl)
	mockDataStore := mocks.NewMockDatastore(mockCtrl)
	env := config.Env{mockDataStore, mockSessionStore}

	u := models.User{}
}
