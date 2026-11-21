package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
)

type contextKey string

const requestIDKey contextKey = "request_id"
const requestIDHeader = "X-Request-ID"

var validRequestID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)

		if !validRequestID.MatchString(requestID) {
			var err error

			requestID, err = generateRequestID()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		w.Header().Set(requestIDHeader, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)

	if !ok {
		return ""
	}
	return requestID
}

func generateRequestID() (string, error) {
	var randomBytes [16]byte
	fmt.Println(randomBytes)

	if _, err := rand.Read(randomBytes[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes[:]), nil
}
