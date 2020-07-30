package main

import (
	"os"

	_ "github.com/giannimassi/shorturl/docs"
	"github.com/giannimassi/shorturl/pkg/routes"
	"github.com/giannimassi/shorturl/pkg/storage"
)

func main() {
	if err := run(); err != nil {
		println("Unexpected error:", err)
		os.Exit(1)
	}
}

func run() error {
	return routes.Start(storage.NewMemoryStore())
}
