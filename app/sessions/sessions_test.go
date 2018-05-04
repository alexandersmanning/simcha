package sessions

import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/alexandersmanning/simcha/app/models"
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

func clearSessions(t *testing.T) {
	req, _ := http.NewRequest("POST", "/logout", nil)
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

func TestLogin(t *testing.T) {
	clearSessions(t)

	req, _ := http.NewRequest("POST", "/login", nil)
	rec := httptest.NewRecorder()

	u := models.User{ID: 100, SessionToken: "fakeToken"}

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
	//rec := httptest.NewRecorder()

	//mockCtrl := gomock.NewController(t)
	//mockDataStore := mocks.NewMockDatastore(mockCtrl)

	t.Run("Logout without without user", func(t *testing.T) {
		clearSessions(t)
		setLoginCredentials(t, req,100, "fakeToken")

		sessions, err := session.Get(req, "session")
		if err != nil {
			t.Fatal(err)
		}

		val := sessions.Values["id"]

		if id, ok := val.(int); !ok {
			t.Errorf("Setup failed")
		} else if id != 100 {
			t.Errorf("Setup failed, expected %d, got %d", 100, id)
		}

		val = sessions.Values["token"]
		if token, ok := val.(string); !ok {
			t.Errorf("Setup failed")
		} else if token != "fakeToken" {
			t.Errorf("Setup failed expeced %s, got %s", "fakeToken", token)
		}

		//session.Logout(&models.User{}, mockDataStore, rec, req)
	})
}
