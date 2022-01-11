package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Mux *mux.Router
}

func NewServer(m *mux.Router) (*Server, error) {
	s := &Server{
		Mux: m,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}
