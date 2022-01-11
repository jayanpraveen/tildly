package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Mux *mux.Router
}

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}
