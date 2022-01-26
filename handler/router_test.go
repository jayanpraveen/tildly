package handler

import (
	"testing"

	"github.com/gorilla/mux"
)

func Test_RunRouter(t *testing.T) {
	t.Run("Run Router", func(t *testing.T) {
		svr := &Server{
			Mux: mux.NewRouter(),
		}
		mx := NewRouter(svr)
		mx.RunRouter()
	})
}
