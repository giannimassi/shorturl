package routes

import (
	"net/http"
	"net/url"
	"strings"
)

// ShortURLProvider is the repository from which short url are fetched.
type ShortURLProvider interface {
	// ShortURL returns true if a short url is found for the provided key, false
	ShortURL(key string) (*url.URL, bool)
}

// RedirectHandler implements a handler that redirects all urls to example.com
// NOTE: only GET requests are supported and tested.
// Reference: https://tools.ietf.org/html/rfc7231#section-6.4.2
func RedirectHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shortURL, found := s.ShortURL(strings.TrimPrefix(r.URL.Path, "/"))
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, r, shortURL.String(), http.StatusMovedPermanently)
	})
}

// Middlewares

// OnlyIf is a middleware that responds with 405 to any method other than the one provided
func OnlyIf(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
