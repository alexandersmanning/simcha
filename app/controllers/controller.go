package controllers

import (
	"fmt"
	"net/http"
)

// Response Object for standardizing JSON output
type JSONResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"result":"", "error": %q}`, err.Error())
}
