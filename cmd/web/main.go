package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"html/template"

	"github.com/JerryLegend254/snippetbox/pkg/models"
	"github.com/JerryLegend254/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {

	addr := flag.String("addr", ":8000", "HTTP network address")
	dsn := flag.String("dsn", "web:Pass123@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	//session.Secure = true

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	//tlsConfig := &tls.Config{
	//	PreferServerCipherSuites: true,
	//	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	//}

	infoLog.Println("Starting server on: ", *addr)

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Use the Ping() method to check if the connection is still alive.
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
