package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"Service/internal/config"
	"Service/internal/logger"
	"Service/internal/models"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog       *log.Logger
	logger         *slog.Logger
	config         *config.Config
	users          *models.UserModel
	days           *models.DayModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	logger := logger.New()
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "../../config.env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}

	conf := config.New()

	addr := flag.String("addr", conf.Server.Addr, "http network address")
	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", conf.BD.User, conf.BD.Pass, conf.BD.Addr, conf.BD.DBName)
	dsn := flag.String("dsn", psqlInfo, "database source name")
	flag.Parse()

	// init pgx
	dbpool, err := NewPG(context.Background(), *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer dbpool.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()

	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		errorLog:       errorLog,
		config:         conf,
		users:          &models.UserModel{Pg: dbpool, Ctx: context.Background()},
		days:           &models.DayModel{Pg: dbpool, Ctx: context.Background()},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//app.Logf(slog.LevelDebug, "Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
