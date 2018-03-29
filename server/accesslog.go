package server

import (
	"net/http"
	"time"
)

// wrapper handler for logging purposes
// abandon due to difficulty retrieiving status code from ResponseWriter
func accessLog(inner http.Handler, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, req)

		log.Infof(
			"%s %s %s %s %s",
			req.RemoteAddr,
			req.Method,
			req.RequestURI,
			name,
			time.Since(start),
		)
	})
}
