package handler

import (
	"net/http"

	"github.com/jayanpraveen/tildly/datastore"
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
	sr := r.PathPrefix("/api").Subrouter()

	// todo
	rd := datastore.DialRedisCache()
	ch := service.NewCacheRepo(rd)
	us := service.NewUrlService(ch)
	uh := NewUrlHandler(us)

	r.HandleFunc("/", uh.handleIndex())
	sr.HandleFunc("/", uh.handleLongUrl()).Queries("longUrl", "{longUrl}")
	r.HandleFunc("/{hash}", uh.handleShortUrl())

	r.Use(middleware.LoggingMiddleware)
	sr.Use(middleware.ValidateUrlMiddleware)

	return http.ListenAndServe(":8080", r)
}
