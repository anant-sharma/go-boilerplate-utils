package server

import (
	"context"
	"net/http"

	"github.com/anant-sharma/go-utils"
)

type key string

const (
	requestIDKey key = "X-Request-Id"
)

func tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(string(requestIDKey))
		if requestID == "" {
			requestID = utils.GenerateShortID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		w.Header().Set(string(requestIDKey), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
