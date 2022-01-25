package handler

import (
	"net/http"

	"github.com/gorilla/mux"
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

func (rtr *router) RunRouter() *mux.Router {
	r := rtr.srv.Mux
	sr := r.PathPrefix("/api").Subrouter()

	rd := datastore.DialRedisClient()
	ch := service.NewCacheRepo(rd)
	us := service.NewUrlService(ch, datastore.NewEtcd())
	uh := NewUrlHandler(us)

	r.HandleFunc("/", uh.handleIndex())
	sr.HandleFunc("/longUrl", uh.handleLongUrl()).Methods(http.MethodPost)
	r.HandleFunc("/{hash}", uh.handleShortUrl())

	r.Use(middleware.LoggingMiddleware)

	return r
}
