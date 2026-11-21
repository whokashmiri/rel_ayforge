package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	var output bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(&output, nil),
	)

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusCreated)
	})

	handler := RequestID(Logging(logger)(next))

	req := httptest.NewRequest(
		http.MethodGet,
		"/test?token=secret",
		nil,
	)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}

	var record map[string]any

	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatalf("decode log JSON: %v", err)
	}

	if record["method"] != "GET" {
		t.Errorf("method = %v, want GET", record["method"])
	}

	if record["path"] != "/test" {
		t.Errorf("path = %v, want /test", record["path"])
	}

	if record["status"] != float64(201) {
		t.Errorf("status = %v, want 201", record["status"])
	}

	requestID := rec.Header().Get("X-Request-ID")

	if record["request_id"] != requestID {
		t.Errorf("logged request ID does not match response header")
	}

	if _, ok := record["duration_ms"]; !ok {
		t.Error("duration_ms missing from log")
	}

	if bytes.Contains(output.Bytes(), []byte("secret")) {
		t.Error("sensitive query value appeared in log")
	}
}
