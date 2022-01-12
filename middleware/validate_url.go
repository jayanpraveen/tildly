package middleware

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func isValidUrl(longUrl string) bool {
	if u, err := url.Parse(longUrl); err == nil && u.Scheme != "" && u.Host != "" {
		return true
	}
	return false
}

func ValidateUrlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if !isValidUrl(vars["longUrl"]) {
			return
		}

		next.ServeHTTP(w, r)

	})
}
