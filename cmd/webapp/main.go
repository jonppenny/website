package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/spf13/viper"
	"html/template"
	"jonppenny.co.uk/webapp/internal/templates"
	"jonppenny.co.uk/webapp/pkg/models/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog           *log.Logger
	infoLog            *log.Logger
	session            *sessions.Session
	posts              *mysql.PostModel
	templateCache      map[string]*template.Template
	adminTemplateCache map[string]*template.Template
	users              *mysql.UserModel
}

func main() {
	addr := flag.String("addr", ":9990", "HTTP network address")
	config := flag.String("config", "config.yaml", "Config file for default settings.")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)

	viper.SetConfigFile(*config)
	err := viper.ReadInConfig()
	if err != nil {
		errorLog.Fatal(err)
	}

	secret := viper.Get("secret").(string)

	dsn := fmt.Sprintf(
		"%s:%s@%s(%s)/%s",
		viper.Get("username"),
		viper.Get("password"),
		viper.Get("protocol"),
		viper.Get("host"),
		viper.Get("database"),
	)
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := templates.NewTemplateCache("./web/html/webapp")
	if err != nil {
		errorLog.Fatal(err)
	}

	adminTemplateCache, err := templates.NewTemplateCache("./web/html/admin")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:           errorLog,
		infoLog:            infoLog,
		session:            session,
		posts:              &mysql.PostModel{DB: db},
		templateCache:      templateCache,
		adminTemplateCache: adminTemplateCache,
		users:              &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
