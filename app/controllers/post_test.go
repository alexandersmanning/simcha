package controllers

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/mocks"
	"github.com/alexandersmanning/simcha/app/models"
)

//func (mdb *mockDB) AllPosts() ([]*models.Post, error) {
//	posts := []*models.Post{}
//	posts = append(posts, &models.Post{Body: "Body Post 1", Title: "Title Post 1"})
//	posts = append(posts, &models.Post{Body: "Body Post 2", Title: "Title Post 2"})
//	return posts, nil
//}
//
//func (mdb *mockDB) CreatePost(p models.Post) error {
//	return nil
//}

func TestPostIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	env := config.Env{DB: mockDatastore}

	posts := []*models.Post{}
	posts = append(posts, &models.Post{Body: "Body Post 1", Title: "Title Post 1"})
	posts = append(posts, &models.Post{Body: "Body Post 2", Title: "Title Post 2"})

	mockDatastore.EXPECT().AllPosts().Return(posts, nil)

	PostIndex(&env)(rec, req, nil)

	if rec.Code != 200 {
		t.Errorf("Expected a status of 200, got %d", rec.Code)
	}

	if val, ok := rec.HeaderMap["Content-Type"]; !ok {
		t.Error("Expected a content type header, got nothing")
	} else if len(val) > 1 || val[0] != "application/json" {
		t.Errorf("Expected content type of %s, got %s", "application/json", val)
	}

	t.Error(rec.Body.String())
}
