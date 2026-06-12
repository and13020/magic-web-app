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

	db, err := setupDB(c.DbConfig.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Note: to store custom types in our cookie, we would register it prior to creating session (in main or init)
	store := sessions.NewCookieStore([]byte(c.SessionAuthKey), []byte(c.SessionEncrKey))
	store.Options = &sessions.Options{
		SameSite: http.SameSiteLaxMode,
		Path:     "/",      // cookie available across domain
		MaxAge:   3600 * 8, // 8 hrs
		Secure:   true,     // set to TRUE if using https
		HttpOnly: true,
	}

	app := &application{
		errorLog:    log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		card:        r.NewCardRepository(db),
		tmplDir:     "./templates",
		publicPath:  "./public/",
		store:       store,
		user:        r.NewUserRepository(db),
		sessionName: c.SessionName,
	}

	app.tp = NewTemplateRenderer(app.tmplDir)

	app.Serve()
}
