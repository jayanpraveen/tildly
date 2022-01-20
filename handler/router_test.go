package handler

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/jayanpraveen/tildly/server"
)

func Test_RunRouter(t *testing.T) {
	t.Run("Run Router", func(t *testing.T) {
		svr := &server.Server{
			Mux: mux.NewRouter(),
		}
		mx := NewRouter(svr)
		mx.RunRouter()
	})
}
