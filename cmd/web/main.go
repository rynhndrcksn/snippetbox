package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// application struct holds the application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP port")
	staticDir := flag.String("staticDir", "./ui/static", "Static assets directory")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes(*staticDir))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
