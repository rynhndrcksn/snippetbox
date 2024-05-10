package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rynhndrcksn/snippetbox/internal/models"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// application struct holds the application-wide dependencies
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP port")
	staticDir := flag.String("staticDir", "./ui/static", "Static assets directory")
	dbConn := flag.String("dbConn", "web:pass@/snippetbox?charset=utf8&parseTime=true", "Database connection string")
	flag.Parse()

	// Initialize a new structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize the database connection
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

	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new decoder instance.
	formDecoder := form.NewDecoder()

	// Initialize a new session manager with MySQL as the session store.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Add everything to our 'application' struct.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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
