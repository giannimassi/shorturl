package main

import (
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		println("Unexpected error:", err)
		os.Exit(1)
	}
}

func run() error {
	http.ListenAndServe(":80", http.HandlerFunc(redirect))
	return nil
}

const redirectTo = "https://example.com/"

// redirect implements a handler that redirects all urls to example.com
// NOTE: only GET requests are supported and tested.
// Reference: https://tools.ietf.org/html/rfc7231#section-6.4.2
func redirect(w http.ResponseWriter, r *http.Request) {
	// Only GET is supported
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.Redirect(w, r, redirectTo, http.StatusMovedPermanently)
}
