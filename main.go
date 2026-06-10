package main

import (
	"log"
	r "magic/repository"
	"net/http"
	"os"

	"magic/config"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	card        *r.CardRepository
	tmplDir     string
	tp          *TemplateRenderer
	publicPath  string
	store       *sessions.CookieStore
	user        *r.UserRepository
	sessionName string
}

func main() {

	c := config.LoadEnv()

	// TODO: access control for DB? for sqlite3 does it matter?
	db, err := setupDB(c.DbConfig.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: add key as env var or secrets
	store := sessions.NewCookieStore([]byte(c.SessionKey))
	store.Options = &sessions.Options{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/", // cookie available across domain
		MaxAge:   86400 * 5,
		// Secure: true, // set to TRUE if using https
	}

	app := &application{
		errorLog:    log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		card:        r.NewCardRepository(db), // Initialize with any dependencies needed for the repository
		tmplDir:     "./templates",
		publicPath:  "./public/",
		store:       store,
		user:        r.NewUserRepository(db),
		sessionName: c.SessionName,
	}

	app.tp = NewTemplateRenderer(app.tmplDir)

	app.Serve()
}
