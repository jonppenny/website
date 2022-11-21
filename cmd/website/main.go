package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/spf13/viper"
	"html/template"
	"jonppenny.co.uk/webapp/internal/database"
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
	errorLog                 *log.Logger
	infoLog                  *log.Logger
	session                  *sessions.Session
	posts                    *mysql.PostModel
	pages                    *mysql.PageModel
	users                    *mysql.UserModel
	menus                    *mysql.MenuModel
	templateCache            map[string]*template.Template
	adminTemplateCache       map[string]*template.Template
	credentialsTemplateCache map[string]*template.Template
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

	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?parseTime=true",
		viper.Get("username"),
		viper.Get("password"),
		viper.Get("host"),
		viper.Get("database"),
	)
	db, err := database.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := templates.NewTemplateCache("./static/html/web")
	if err != nil {
		errorLog.Fatal(err)
	}

	adminTemplateCache, err := templates.NewTemplateCache("./static/html/admin")
	if err != nil {
		errorLog.Fatal(err)
	}

	credentialsTemplateCache, err := templates.NewTemplateCache("./static/html/credentials")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(viper.Get("secret").(string)))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:                 errorLog,
		infoLog:                  infoLog,
		session:                  session,
		posts:                    &mysql.PostModel{DB: db},
		pages:                    &mysql.PageModel{DB: db},
		users:                    &mysql.UserModel{DB: db},
		templateCache:            templateCache,
		adminTemplateCache:       adminTemplateCache,
		credentialsTemplateCache: credentialsTemplateCache,
	}

	tlsConfig := &tls.Config{
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
