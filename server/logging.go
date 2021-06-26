package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func logging(next http.Handler) http.Handler {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		defer func() {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}
			log.Println(requestID, r.Method, r.URL.Path, time.Since(startTime))
		}()
		next.ServeHTTP(w, r)
	})
}
