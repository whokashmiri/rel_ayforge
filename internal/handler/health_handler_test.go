package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet, "/health", nil,
	)
	recorder := httptest.NewRecorder()
	Health(recorder, req)
	assertResponse(t, recorder, "ok")
}

func TestReady(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet, "/ready", nil,
	)
	recorder := httptest.NewRecorder()
	Ready(recorder, req)
	assertResponse(t, recorder, "ready")
}

func assertResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus string) {
	t.Helper()

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("status code = %d , want 200", result.StatusCode)
	}
	contentType := result.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("content type = %q, want application/json", contentType)
	}
	var response HealthResponse
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("decode JSON response : %v", err)
	}
	if response.Status != expectedStatus {
		t.Errorf("response status = %q , want %q", response.Status, expectedStatus)
	}
}
