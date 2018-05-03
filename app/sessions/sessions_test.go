package sessions

import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
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

func TestLogin(t *testing.T) {
	req, _ := http.NewRequest("POST", "/login", nil)
	rec := httptest.NewRecorder()
}
