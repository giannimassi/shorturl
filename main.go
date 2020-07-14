package main

import (
	"net/http"
	"net/url"
	"os"

	"github.com/giannimassi/shorturl/pkg/routes"
)

func main() {
	if err := run(); err != nil {
		println("Unexpected error:", err)
		os.Exit(1)
	}
}

func run() error {
	return http.ListenAndServe(":80", routes.OnlyIf("GET", routes.RedirectHandler(&acceptAll{})))
}

const redirectTo = "https://example.com/"

//acceptAll is a simple mock providing the short-url for the redirect handler
type acceptAll struct{}

func (s *acceptAll) ShortURL(key string) (*url.URL, bool) {
	if key == "err" {
		return nil, false
	}
	url, err := url.Parse(redirectTo)
	if err != nil {
		panic(err)
	}
	return url, true
}
