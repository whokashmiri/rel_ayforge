package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRequestIDGenerated(t *testing.T) {
	var contextID string

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		contextID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	RequestID(next).ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	headerID := result.Header.Get("X-Request-ID")

	if !regexp.MustCompile(`^[a-f0-9]{32}$`).MatchString(headerID) {
		t.Fatalf("invalid generated request ID: %q", headerID)
	}

	if contextID != headerID {
		t.Errorf("context ID %q != header ID %q", contextID, headerID)
	}
}

func TestRequestIDAccepted(t *testing.T) {
	const expectedID = "client-request-123"

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if got := GetRequestID(r.Context()); got != expectedID {
			t.Errorf("context request ID = %q", got)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-ID", expectedID)

	rec := httptest.NewRecorder()

	RequestID(next).ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != expectedID {
		t.Errorf("response request ID = %q, want %q", got, expectedID)
	}
}

func TestInvalidRequestIDReplaced(t *testing.T) {
	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-ID", "invalid request id")

	rec := httptest.NewRecorder()

	RequestID(next).ServeHTTP(rec, req)

	got := rec.Header().Get("X-Request-ID")

	if got == "" || strings.Contains(got, " ") {
		t.Errorf("invalid generated replacement ID: %q", got)
	}

	if got == "invalid request id" {
		t.Error("invalid client request ID was accepted")
	}
}
