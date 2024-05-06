package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// setup flags and parse them
	addr := flag.String("addr", ":8080", "HTTP port")
	staticDir := flag.String("staticDir", "./ui/static", "Static assets directory")
	flag.Parse()

	// setup slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory
	// path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir(*staticDir))

	// use mux.Handle() to register the file server as a handler for all URL paths that start with "/static/"
	// for match paths, we strip the "/static" prefix before the request reaches the file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register application routes
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
