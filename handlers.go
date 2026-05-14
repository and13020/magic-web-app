package main

import (
	"fmt"
	"magic/repository"
	"magic/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// mux.HandleFunc("/", app.home)
// mux.HandleFunc("/search/{name}", app.card.search) // {name} is a path variable that can be accessed in the handler function
// mux.HandleFunc("/random", app.card.random)
// mux.HandleFunc("/update", app.card.updateCard)
// mux.HandleFunc("/delete", app.card.deleteCard)

const loggedInUserKey = "auth"

// GIN handlers require *gin.Context, giving methods for request/response
func (app *application) home(c *gin.Context) {
	app.render(c, "index.html", nil)
}

func (app *application) random(c *gin.Context) {
	card, err := app.card.GetRandomCard()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching random card: %v", err)
		return
	}

	err = app.card.SaveCard(card)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching random card: %v", err)
		return
	}

	app.render(c, "card.html", &templateData{Card: card})
}

func (app *application) getCards(c *gin.Context) {

	// call GET CARDS
	// app.render (templateData) will have the cards
	// call proper html file which can display multiple cards

	//TODO: get name from form field or from the api itself

	app.render(c, "index.html", nil)

}

// getCardsForm is called for POST requests on "/search"
// It reads form input and displays the data back to the user
func (app *application) getCardsForm(c *gin.Context) {
	// name := "black lotus"
	// cards, err := app.card.GetCardsByName(name)
	// if err != nil {
	// 	c.String(http.StatusBadRequest, "Error fetching card by name: ", name, " --- ", err)
	// }

	form := NewForm(c.Request.Form)
	// validation on form fields..
	// ie fewer than 1000 chars for search

	// once validation attempted, check if any errors found
	if !form.Valid() {
	}

	err := c.Request.ParseForm()
	if err != nil {
		fmt.Println("Could not parse form! ", err)
		form.Errors.Add("generic", "could not parse form")

		app.render(c, "index.html", &templateData{Form: form})
		return
	}

	// if form is valid w/ no errors
	// read from it
	// use data

	name := c.Request.FormValue("name")
	cards, err := app.card.GetCardsByName(name)
	if err != nil {
		app.render(c, "index.html", &templateData{Cards: []repository.Card{}})
	}

	fmt.Println("Form has name: ", name)
	fmt.Println("Cards list: ", cards)
	app.render(c, "index.html", &templateData{Cards: cards})
}

func (app *application) Signup(c *gin.Context) {

	if c.Request.Method == http.MethodPost {

		err := c.Request.ParseForm()
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		form := NewForm(c.Request.Form)

		// validations
		form.Required("email", "password", "password2", "username").
			MinLength("password", 8).
			MaxLength("password", 255).
			MinLength("password2", 8).
			MaxLength("password2", 255).
			MatchPass("password", "password2").
			MinLength("username", 3)

		if !form.Valid() {
			app.render(c, "signup.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		email := c.Request.FormValue("email")
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")

		// Check if email or username already used
		err = app.user.Validate(email, username)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		// Add user to DB
		err = app.user.Add(email, password, username)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		// 303 redirect (indicating POST to GET)
		http.Redirect(c.Writer, c.Request, "/login", http.StatusSeeOther)
		return

	}

	app.render(c, "signup.html", nil)
}

func (app *application) Login(c *gin.Context) {

	// If session exists, redirect to home
	if app.isAuthenticated(c.Request) {
		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
		return
	}

	if c.Request.Method == http.MethodPost {

		err := c.Request.ParseForm()
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			app.errorLog.Printf("Could not parse form: %v", err)
			return
		}

		form := NewForm(c.Request.Form)
		form.Required("email", "password")

		if !form.Valid() {
			app.errorLog.Printf("Validation failed: %v", form.Errors)
			app.render(c, "login.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		email := c.Request.FormValue("email")
		password := c.Request.FormValue("password")

		// check if user exists in DB
		user, err := app.user.GetUserByField("email", email)
		if err != nil {
			app.errorLog.Printf("Login failed: %v", err)
			app.render(c, "login.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		// check if pass from DB/input match
		if !utils.CheckPassword(user.Password, password) {
			app.errorLog.Printf("Login failed due to incorrect password")
			app.render(c, "login.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		// Create session
		session, err := app.store.Get(c.Request, session_key)
		if err != nil {
			fmt.Println("Failed to store.Get to create new session: ", err)
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values[loggedInUserKey] = true

		err = session.Save(c.Request, c.Writer)
		if err != nil {
			fmt.Println("session failed to save: ", err)
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Session created: ", session)
		app.SetFlash(c, "Successfully logged in")

		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
		return
	}

	app.render(c, "login.html", nil)

}

// Logoff deletes session and redirects to login page
func (app *application) Logoff(c *gin.Context) {
	s, _ := app.store.Get(c.Request, session_key)
	s.Options.MaxAge = -1

	fmt.Println("Successfully logged off")
	s.Save(c.Request, c.Writer)

	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
	// app.render("put in data for flash in template")
}
