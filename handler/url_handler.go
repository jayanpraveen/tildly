package handler

import (
	"fmt"
	"net/http"

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
