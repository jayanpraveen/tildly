package main

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/jayanpraveen/tildly/handler"
	"github.com/jayanpraveen/tildly/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	srv, err := server.NewServer(mux.NewRouter())
	if err != nil {
		return err
	}

	rtr := handler.NewRouter(srv)

	return rtr.RunRouter()
}
