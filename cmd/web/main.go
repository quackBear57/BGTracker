package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"bgtracker.ndr54.com/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	debug				bool
	logger				*slog.Logger
	games				*models.GameModel
	users				*models.UserModel
	templateCache		map[string]*template.Template
	formDecoder			*form.Decoder
	sessionManager		*scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network port")
	dsn := flag.String("dsn",
	"bgtracker_ur:spree-sculpture-commuting@tcp(192.168.0.3:3306)/bgtracker?parseTime=true",
	"data source name")
	debug := flag.Bool("debug", false, "Enable debug logs")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		debug:				*debug,
		logger:				logger,
		games:				&models.GameModel{DB: db},
		users:				&models.UserModel{DB: db},
		templateCache:		templateCache,
		formDecoder:		formDecoder,
		sessionManager:		sessionManager,
	}

	srv := &http.Server{
		Addr:			*addr,
		Handler:		app.routes(),
		ErrorLog:		slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:	time.Minute,
		ReadTimeout:	time.Second * 5,
		WriteTimeout:	time.Second * 10,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}