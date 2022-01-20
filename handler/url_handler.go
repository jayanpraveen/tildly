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
	service.UrlRepository
}

// note: using `service.UrlRepository` instead of directly unsing `service.UrlService` to make tests/mocks easier
func NewUrlHandler(us service.UrlRepository) *UrlHandler {
	return &UrlHandler{
		UrlRepository: us,
	}
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
			panic(err)
		}

		if err := isValidUrl(u.LongUrl); err != nil {
			panic(err)
		}

		if err := s.SaveUrl(u.LongUrl); err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Url created!")
	}
}

func (s *UrlHandler) handleShortUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		if vars["hash"] == "api" {
			return
		}

		u, err := s.GetUrlByHash(vars["hash"])
		if err != nil {
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
