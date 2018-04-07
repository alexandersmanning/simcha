package controllers

import (
	"fmt"
	"net/http"
)

func badReq(w http.ResponseWriter, isJSON bool, err error, status int) {
	if !isJSON {
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"result":"", "error": %q}`, err.Error())
}
