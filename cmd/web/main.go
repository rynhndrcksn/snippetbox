package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/rynhndrcksn/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// application struct holds the application-wide dependencies
type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP port")
	staticDir := flag.String("staticDir", "./ui/static", "Static assets directory")
	dbConn := flag.String("dbConn", "web:pass@/snippetbox?charset=utf8&parseTime=true", "Database connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dbConn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes(*staticDir))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
