package routes

import "net/http"

const redirectTo = "https://example.com/"

// RedirectHandler implements a handler that redirects all urls to example.com
// NOTE: only GET requests are supported and tested.
// Reference: https://tools.ietf.org/html/rfc7231#section-6.4.2
func RedirectHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTo, http.StatusMovedPermanently)
	})
}

// Middlewares

// AllowGETOnly is a middleware that responds with 405 to any method other than GET
func AllowGETOnly(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
