//Entry point for controller tests, contains helper methods
package controllers

import (
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
