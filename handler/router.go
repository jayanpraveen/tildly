package handler

import (
	"net/http"

	"github.com/jayanpraveen/tildly/middleware"
	"github.com/jayanpraveen/tildly/server"
	"github.com/jayanpraveen/tildly/service"
)

type router struct {
	srv *server.Server
}

func NewRouter(s *server.Server) *router {
	return &router{
		srv: s,
	}
}

func (rtr *router) RunRouter() error {
	r := rtr.srv.Mux

	us := service.NewUrlService()
	uh := NewUrlHandler(us)

	r.HandleFunc("/", uh.handleIndex())
	// ...

	r.Use(middleware.LoggingMiddleware)

	return http.ListenAndServe(":8080", r)
}
