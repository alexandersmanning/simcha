package sessions

import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/alexandersmanning/simcha/app/models"
	"github.com/golang/mock/gomock"
	"github.com/alexandersmanning/simcha/app/mocks"
)

func TestMain(m *testing.M) {
	setupSessions()
	exitCode := m.Run()

	os.Exit(exitCode)
}

var session *Session

func setupSessions() {
	session = InitStore(os.Getenv("APPLICATION_SECRET"))
}

func clearSessions(t *testing.T, req *http.Request) {
	sessions, err := session.Get(req, "session")

	if err != nil {
		t.Fatal(err)
	}

	sessions.Values["id"] = nil;
	sessions.Values["token"] = nil;
}

func setLoginCredentials(t *testing.T, req *http.Request, id int, token string) {
	sessions, err := session.Get(req, "session")

	if err != nil {
		t.Fatal(err)
	}

	sessions.Values["id"] = id
	sessions.Values["token"] = token
}

func verifySetLoginCredentials(t *testing.T, req *http.Request, id int, token string) {
	t.Helper()

	sessions, err := session.Get(req, "session")
	if err != nil {
		t.Fatal(err)
	}

	val := sessions.Values["id"]

	if foundId, ok := val.(int); !ok {
		t.Errorf("Setup failed")
	} else if foundId != id {
		t.Errorf("Setup failed, expected %d, got %d", id, foundId)
	}

	val = sessions.Values["token"]
	if foundToken, ok := val.(string); !ok {
		t.Errorf("Setup failed")
	} else if foundToken != token {
		t.Errorf("Setup failed expeced %s, got %s", token, foundToken)
	}
}

func TestLogin(t *testing.T) {

	req, _ := http.NewRequest("POST", "/login", nil)
	rec := httptest.NewRecorder()

	u := models.User{ID: 100, SessionToken: "fakeToken"}

	clearSessions(t, req)

	if err := session.Login(&u, rec, req); err != nil {
		t.Fatal(err)
	}

	testSession, err := session.Get(req, "session")

	if err != nil {
		t.Fatal(err)
	}

	val := testSession.Values["id"]

	if id, ok := val.(int); !ok {
		t.Errorf("Expected a val of %d, found nothing", u.ID)
	} else if id != u.ID {
		t.Errorf("Expected a val of %d, got %d", u.ID, id)
	}

	val = testSession.Values["token"]

	if token, ok := val.(string); !ok {
		t.Errorf("Expected a val of %s, found nothing", u.SessionToken)
	} else if token != u.SessionToken {
		t.Errorf("Expected a val of %s, got %s", u.SessionToken, token)
	}

}

func TestLogout(t *testing.T) {
	req, _ := http.NewRequest("GET", "/logout", nil)
	rec := httptest.NewRecorder()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDataStore := mocks.NewMockDatastore(mockCtrl)

	t.Run("Logout with logged in user", func(t *testing.T) {
		u := &models.User{ID: 100, SessionToken: "fakeToken" }
		clearSessions(t, req)
		setLoginCredentials(t, req, u.ID, u.SessionToken)

		verifySetLoginCredentials(t, req, u.ID, u.SessionToken)

		mockDataStore.EXPECT().UpdateSessionToken(u.ID).Return(nil).Times(1)
		session.Logout(u, mockDataStore, rec, req)

		sessions, err := session.Get(req, "session")

		if err != nil {
			t.Fatal(err)
		}

		if val := sessions.Values["id"]; val != nil {
			t.Errorf("Expected ID to be nil after logout, got %v", val)
		}

		if val := sessions.Values["token"]; val != nil {
			t.Errorf("Expected token to be nil after logout, got %v", val)
		}
	})

	t.Run("Without logged in user, session is still set to 0", func(t *testing.T) {
		u := &models.User{ID: 100, SessionToken: "fakeToken"}
		clearSessions(t, req)
		setLoginCredentials(t, req,u.ID, u.SessionToken)

		verifySetLoginCredentials(t, req, u.ID, u.SessionToken)

		mockDataStore.EXPECT().UpdateSessionToken(u.ID).Return(nil).Times(0)
		session.Logout(&models.User{}, mockDataStore, rec, req)

		sessions, err := session.Get(req, "session")

		if err != nil {
			t.Fatal(err)
		}

		if val := sessions.Values["id"]; val != nil {
			t.Errorf("Expected ID to be nil after logout, got %v", val)
		}

		if val := sessions.Values["token"]; val != nil {
			t.Errorf("Expected token to be nil after logout, got %v", val)
		}
	})
}
