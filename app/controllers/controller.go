package controllers

import (
	"fmt"
	"net/http"
)

func jsonError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"result":"", "error": %q}`, err.Error())
}
