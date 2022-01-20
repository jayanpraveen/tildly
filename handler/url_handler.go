package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jayanpraveen/tildly/service"
)

type UrlHandler struct {
	urs service.UrlRepository
}

func NewUrlHandler(us service.UrlRepository) *UrlHandler {
	return &UrlHandler{urs: us}
}

func isValidUrl(longUrl string) error {
	if longUrl == "" {
		return fmt.Errorf("")
	}

	if u, err := url.Parse(longUrl); err == nil && u.Scheme != "" && u.Host != "" {
		return nil
	}
	return fmt.Errorf("")
}

func (s *UrlHandler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "tildly !")
	}
}

func (s *UrlHandler) handleLongUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type PostUrl struct {
			LongUrl string `json:"longUrl"`
		}

		var u PostUrl

		if err := service.DecodeJson(&u, r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.SetError(http.StatusBadRequest, "Not a proper JSON format", w)
			return
		}

		if err := isValidUrl(u.LongUrl); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.SetError(http.StatusBadRequest, "Not a valid URL", w)
			return
		}

		if err := s.urs.SaveUrl(u.LongUrl); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			service.SetError(http.StatusInternalServerError, "Internal server error", w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "tildly url created!")
	}
}

func (s *UrlHandler) handleShortUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		if vars["hash"] == "api" {
			return
		}

		u, err := s.urs.GetUrlByHash(vars["hash"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			notFoundTemplate(w, r)
			return
		}

		http.Redirect(w, r, u.LongUrl, http.StatusPermanentRedirect)
	}
}

func notFoundTemplate(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./handler/NotFound.html"))
	t.Execute(w, nil)
}
