package api

import (
	"encoding/json"
	"net/http"
)

// RegisterRoutes wires API handlers onto the provided mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/api/v1/message", messageHandler)
	mux.HandleFunc("/api/v1/toys", toysCollectionHandler)
	mux.HandleFunc("/api/v1/toys/", toysItemHandler)
}

type healthResponse struct {
	Status string `json:"status"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, messageResponse{Message: "Welcome to the toy API"})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, errorResponse{Error: message})
}
