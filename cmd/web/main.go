package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
	"strconv"
	"fmt"

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
	players				*models.PlayerModel
	sessions			*models.SessionModel
	scoresheets			*models.ScoresheetModel
	templateCache		map[string]*template.Template
	formDecoder			*form.Decoder
	sessionManager		*scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network port")
	
	database_user := os.Getenv("DATABASE_USER")
	database_password := os.Getenv("DATABASE_PASSWORD")
	database_url := os.Getenv("DATABASE_URL")
	database_port := os.Getenv("DATABASE_PORT")
	database_name := os.Getenv("DATABASE_NAME")

	// defaults for testing
	if database_user == "" {database_user = "bgtracker_ur"}
	if database_password == "" {database_password = "spree-sculpture-commuting"}
	if database_url == "" {database_url = "192.168.0.3"}
	if database_port == "" {database_port = "3306"}
	if database_name == "" {database_name = "bgtracker"}

	dsn := flag.String("dsn",
	database_user + ":" + database_password + 
	"@tcp(" + database_url + ":" + database_port + ")/"  + database_name +
	"?parseTime=true",
	"data source name")

	// debug_env, err := strconv.ParseBool(os.Getenv("DEBUG"))
	debug_env, err := strconv.ParseBool("true")
	debug := flag.Bool("debug", debug_env, "Enable debug logs")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Print out environment flags:
	fmt.Printf("Db user: %s\n", database_user)
	// fmt.Printf("Db password: %s\n", "not shown")
	fmt.Printf("Db URL: %s\n", database_url)
	fmt.Printf("Db port: %s\n", database_port)
	fmt.Printf("Db name: %s\n", database_name)
	fmt.Printf("Debug env: %t\n", debug_env)


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
		players:			&models.PlayerModel{DB: db},
		sessions:			&models.SessionModel{DB: db},
		scoresheets:		&models.ScoresheetModel{DB: db},
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