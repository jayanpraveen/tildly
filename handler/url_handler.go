package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jayanpraveen/tildly/service"
)

type UrlHandler struct {
	*service.UrlService
}

func NewUrlHandler(us *service.UrlService) *UrlHandler {
	return &UrlHandler{UrlService: us}
}

func (s *UrlHandler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "tildly !")
	}
}

func (s *UrlHandler) handleLongUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		err := s.UrlService.SaveUrl(vars["longUrl"])
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(w, "Handle Long Url", vars)
	}
}

func (s *UrlHandler) handleShortUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		if vars["hash"] == "api" {
			return
		}

		u, err := s.UrlService.GetUrlByHash(vars["hash"])
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(w, "Handle Short Url", u.LongUrl)
	}
}
