package main

import (
	"flag"
	"fmt"
	"net/http"
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

	port := flag.String("port", "8080", "http port")
	flag.Parse()

	srv, err := server.NewServer(mux.NewRouter())
	if err != nil {
		return err
	}

	rtr := handler.NewRouter(srv)

	router := rtr.RunRouter()

	return http.ListenAndServe(fmt.Sprintf(":%s", *port), router)
}
