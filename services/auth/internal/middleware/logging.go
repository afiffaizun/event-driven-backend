package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		reqID := GetRequestID(r.Context())
		log.Printf(
			"method=%s path=%s request_id=%s duration=%s",
			r.Method,
			r.URL.Path,
			reqID,
			time.Since(start),
		)
	})
}
