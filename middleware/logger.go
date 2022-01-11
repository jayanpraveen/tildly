package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData.size += size
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.responseData.status = statusCode
}

// The logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{status: 0, size: 0}
		lrw := loggingResponseWriter{ResponseWriter: rw, responseData: responseData}
		next.ServeHTTP(&lrw, r)

		duration := time.Since(start)

		log.SetPrefix("[INFO]: ")
		log.Printf("\nURI: %v, \nDuration: %v, \nStatus: %v, \nSize: %v\n\n",
			r.RequestURI, duration, responseData.status, responseData.size)
	})
}
