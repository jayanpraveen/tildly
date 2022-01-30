package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jayanpraveen/tildly/datastore"
	"github.com/jayanpraveen/tildly/middleware"
	"github.com/jayanpraveen/tildly/service"
)

type router struct {
	srv *Server
}

func NewRouter(s *Server) *router {
	return &router{
		srv: s,
	}
}

func (rtr *router) RunRouter() *mux.Router {
	r := rtr.srv.Mux
	sr := r.PathPrefix("/api").Subrouter()

	etcd := datastore.NewEtcd()
	ch := service.NewCacheRepo(datastore.DialRedisClient())
	csdra := service.NewCassandra(datastore.DialCassandra())
	us := service.NewUrlService(*csdra, *ch, etcd)
	uh := NewUrlHandler(us)

	r.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, us.AC.DisplayCurrentRange())
	})

	r.HandleFunc("/", uh.handleIndex())
	sr.HandleFunc("/longUrl", uh.handleLongUrl()).Methods(http.MethodPost)
	r.HandleFunc("/{hash}", uh.handleShortUrl())

	r.Use(middleware.LoggingMiddleware)

	return r
}
