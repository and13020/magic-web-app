package main

import (
	"fmt"
	"log"
	r "magic/repository"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	card       *r.CardRepository
	tmplDir    string
	tp         *TemplateRenderer
	publicPath string
	store      *sessions.CookieStore
	user       *r.UserRepository
}

var session_key, sessionSecret string

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	session_key = os.Getenv("SESSION_KEY")
	if session_key == "" {
		log.Fatal("SESSION_KEY is required")
	}
	sessionSecret = os.Getenv("SECRET_KEY")
	if sessionSecret == "" {
		log.Fatal("SECRET_KEY is required")
	}
}

func main() {

	loadEnv()

	db, err := setupDB("mtg.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: add key as env var or secrets
	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/", // cookie available across domain
		MaxAge:   86400 * 5,
		// Secure: true, // set to TRUE if using https
	}

	app := &application{
		errorLog:   log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		infoLog:    log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		card:       r.NewCardRepository(db), // Initialize with any dependencies needed for the repository
		tmplDir:    "./templates",
		publicPath: "./public/",
		store:      store,
		user:       r.NewUserRepository(db),
	}

	app.tp = NewTemplateRenderer(app.tmplDir)

	app.Serve()
}

func (app *application) runServer() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	// Register a handler function for the root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	// Start the server on port 8080
	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
