package main

import (
	"fmt"
	"log"
	r "magic/repository"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	card       *r.CardRepository
	tmplDir    string
	tp         *TemplateRenderer
	publicPath string
	session    *sessions.Session
}

func main() {
	// DB WORK
	// 1. create sqlite database
	// 2. create table for cards
	// 3. create function to insert card into database
	// 4. search by any card detail (flexible query)
	// 5. delete card
	// 6. update card
	// 7. data limitations (e.g. only store 1000 cards, delete oldest card when limit is reached)

	// API WORK
	// 1. Created basic request to get card
	// 2. Create other means to search for card
	// 3. Search multiple cards
	// 4. Pagination
	// 5. Cache - check DB if card exists prior to making API call (exact match only)

	// FRONT END
	// 1. Create basic UI to display card details
	// 2. Create search bar to search for card
	// 3. Display multiple cards
	// 4. Pagination
	// 5. Add ability to save card to database (if not already saved)
	// 6. Add ability to delete card from database
	// 7. Save a deck of cards to database (if not already saved)
	// 8. Add ability to delete deck from database
	// 9. Add ability to create deck
	// 10. Add cookies for our sessions to track user data (e.g. saved cards, saved decks, etc)
	// 10. etc etc

	// SIMPLIFY ABOVE:
	// 1. Web server
	// 2. Create API endpoints for above functionality
	// 3. Create front end to call API endpoints and display data.. for now lets just spit out data

	// TODO:
	//Typical Storage Locations
	//Project Root: For simple development, keep the .db file in the same folder as your code.
	//User Data Directories: For production apps, use system-standard folders to ensure the database persists and has the correct permissions:
	//Windows: %AppData%\YourAppName\
	//Linux/macOS: ~/.config/YourAppName/ or ~/.local/share/YourAppName/
	//Mobile Apps (Android/iOS): Stored in the app's private data folder (e.g., /data/data/<package_name>/databases/) to keep it hidden from users.
	db, err := setupDB("mtg.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:   log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		infoLog:    log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		card:       r.NewCardRepository(db), // Initialize with any dependencies needed for the repository
		tmplDir:    "./templates",
		publicPath: "./public/",
		session:    sessions.New([]byte("secret-session")),
	}

	app.tp = NewTemplateRenderer(app.tmplDir)
	app.session.Lifetime = 12 * time.Hour

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
