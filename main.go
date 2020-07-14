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
	http.ListenAndServe(":80", redirectHandler())
	return nil
}
