package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestServer(t *testing.T) {

	s, _ := NewServer(mux.NewRouter())

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	s.ServeHTTP(res, req)

}
