package controllers

import ( "bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/alexandersmanning/simcha/app/config"
	"github.com/alexandersmanning/simcha/app/mocks"
	"github.com/alexandersmanning/simcha/app/models"
)

func TestPostIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	env := config.Env{DB: mockDatastore}

	var posts []*models.Post
	posts = append(posts, &models.Post{Body: "Body Post 1", Title: "Title Post 1"})
	posts = append(posts, &models.Post{Body: "Body Post 2", Title: "Title Post 2"})

	mockDatastore.EXPECT().AllPosts().Return(posts, nil)

	PostIndex(&env)(rec, req, nil)

	checkStatus(rec.Code, 200, t)

	checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

	returnedPosts := []*models.Post{}

	msg, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(msg, &returnedPosts); err != nil {
		t.Fatal(err)
	}

	/*
		this works because of deep equal, however these are two separate pieces of memory
		and therefore would fail any `==` tests
	*/
	if !reflect.DeepEqual(returnedPosts, posts) {
		t.Errorf("Expected %v to equal %v", returnedPosts, posts)
	}
}

func TestPostCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	env := config.Env{DB: mockDatastore}

	post := models.Post{Body: "Test Create Body", Title: "Test Create Title"}
	postJSON, err := json.Marshal(post)
	postBuff := bytes.NewBuffer(postJSON)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", postBuff)

	mockDatastore.EXPECT().CreatePost(post).Return(nil)

	PostCreate(&env)(rec, req, nil)

	checkStatus(rec.Code, 200, t)

	checkHeader(rec.HeaderMap, "Content-Type", "application/json", t)

	msg, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(msg) != `{"result": "success"}` {
		t.Errorf("Expected a successful result, got %q", string(msg))
	}
}
