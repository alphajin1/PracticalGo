package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Copy from chap06 - SurveMux 객체를 Wrapping
func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(w, r)
			d := map[string]interface{}{
				"path":     r.URL.Path,
				"method":   r.Method,
				"bodySize": r.ContentLength,
				"protocol": r.Proto,
				"duration": time.Now().Sub(startTime).Seconds(),
			}

			dJson, _ := json.Marshal(d)

			log.Printf(string(dJson))
		})
}
