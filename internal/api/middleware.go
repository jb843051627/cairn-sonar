package api

import (
	"log"
	"net/http"
	"time"
)

func RequestLog(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf("method=%s path=%s elapsed=%s", r.Method, r.URL.Path, time.Since(started))
	})
}

func Recover(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if value := recover(); value != nil {
				logger.Printf("panic recovered: %v", value)
				Error(w, http.StatusInternalServerError, "panic", "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
