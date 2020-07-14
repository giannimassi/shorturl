package main

import "net/http"

const redirectTo = "https://example.com/"

// redirectHandler implements a handler that redirects all urls to example.com
// NOTE: only GET requests are supported and tested.
// Reference: https://tools.ietf.org/html/rfc7231#section-6.4.2
func redirectHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only GET is supported
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		http.Redirect(w, r, redirectTo, http.StatusMovedPermanently)
	})
}
