//Entry point for controller tests, contains helper methods
package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkStatus(code, expected int, t *testing.T) {
	t.Helper()

	if code != expected {
		t.Errorf("Expected to get a status of %d, got %d instead", expected, code)
	}
}

func checkHeader(headerMap map[string][]string, headerType, expected string, t *testing.T) {
	t.Helper()

	if val, ok := headerMap[headerType]; !ok {
		t.Errorf("Expected header to have %s key, has the following structure: %v", headerType, headerMap)
	} else if val[0] != expected {
		t.Errorf("Expected header key %s to be %s, got %s", headerType, expected, val[0])
	}
}

func TestSendJsonResponse(t *testing.T) {
	t.Run("It sets the CSRF token and Content Type", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/any-url", nil)

		if err != nil {
			t.Fatal(err)
		}

		sendJsonResponse(res, req, []byte("testBody"))

		if v := res.Header().Get("Content-Type"); v != "application/json" {
			t.Errorf("Expected application/json, got %s", v)
		}

		var any bool
		for k := range res.Header() {
			if k == "X-Csrf-Token" {
				any = true
			}
		}
		if !any {
			t.Errorf("Expected a csrf header")
		}

		if v := res.Code; v != 200 {
			t.Errorf("Expected to get %d, got %d", 200, v)
		}
	})
}
