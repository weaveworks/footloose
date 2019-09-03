package api

import (
	"encoding/json"
	"net/http"
)

func sendOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// ErrorResponse is the response API entry points return when they encountered an error.
type ErrorResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{
		Error: err.Error(),
	}
	_ = json.NewEncoder(w).Encode(&resp)
}

// CreatedResponse is the response POST entry points return when a resource has been
// successfully created.
type CreatedResponse struct {
	URI string `json:"uri"`
}

func sendCreated(w http.ResponseWriter, URI string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreatedResponse{
		URI: URI,
	}
	_ = json.NewEncoder(w).Encode(&resp)
}
