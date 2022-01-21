package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(*testing.T) {

	logger_handler := func(w http.ResponseWriter, r *http.Request) {}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	logger_handler(res, req)

	lrw := loggingResponseWriter{
		ResponseWriter: res,
		responseData: &responseData{
			status: 2,
			size:   34,
		},
	}
	lrw.Write([]byte(res.Body.Bytes()))
	lrw.WriteHeader(lrw.responseData.size)

	lw := LoggingMiddleware(http.HandlerFunc(logger_handler))
	lw.ServeHTTP(res, req)

}
