package main

import (
	"net/http"
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
	return http.ListenAndServe(":80", routes.AllowGETOnly(routes.RedirectHandler()))
}
