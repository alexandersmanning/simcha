package controllers

import (
	"encoding/json"
	"net/http"
)

// Response Object for standardizing JSON output
type JSONResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, err error, status int) {
	res := JSONResponse{Error: err.Error() }
	resJSON, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resJSON)
}

func jsonResponse(w http.ResponseWriter, body string) {
	res := JSONResponse{Result: body}
	resJSON, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(resJSON)
}
