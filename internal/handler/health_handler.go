package handler

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, HealthResponse{
		Status: "ok",
	})
}
func Ready(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, HealthResponse{
		Status: "ready",
	})
}

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
