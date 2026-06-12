package main

import (
	"fmt"
	"magic/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserIdKey          = "id"
	formFieldEmail     = "email"
	formFieldPassword  = "password"
	formFieldUsername  = "username"
	formFieldPassword2 = "password2"
)

// GIN handlers require *gin.Context, giving methods for request/response
func (app *application) home(c *gin.Context) {
	app.render(c, "index.html", nil)
}

// random returns a random card to render
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

// TODO: not implemented
// was originally to separate GET/POST implementation.. is it worth it?
// func (app *application) getCards(c *gin.Context) {

// 	// call GET CARDS
// 	// app.render (templateData) will have the cards
// 	// call proper html file which can display multiple cards

// 	//TODO: get name from form field or from the api itself

// 	app.render(c, "index.html", nil)

// }

// getCards is called for GET and POST requests on "/search"
// It reads form input and displays the data back to the user
func (app *application) getCards(c *gin.Context) {

	if c.Request.Method == http.MethodPost {
		err := c.Request.ParseForm()
		if err != nil {
			form := NewForm(c.Request.Form)
			form.Errors.Add("generic", "could not parse form")

			app.render(c, "index.html", &templateData{Form: form})
			return
		}

		form := NewForm(c.Request.Form)
		// validation on form fields..
		// ie fewer than 1000 chars for search

		// once validation attempted, check if any errors found
		if !form.Valid() {
		}

		name := c.Request.FormValue("name")
		cards, err := app.card.GetCardsByName(name)
		if err != nil {
			app.render(c, "index.html", &templateData{Cards: []repository.Card{}})
			return
		}

		app.render(c, "index.html", &templateData{Cards: cards})
		return
	}

	app.render(c, "index.html", nil)
}

func (app *application) signup(c *gin.Context) {

	if c.Request.Method == http.MethodPost {

		err := c.Request.ParseForm()
		if err != nil {
			// Don't expose too much detail to user
			app.errorLog.Println("Failed to parse form during signup: ", err.Error())
			app.render(c, "signup.html", &templateData{Flash: "Failed to read the form 😔"})
			return
		}

		form := NewForm(c.Request.Form)

		// validations
		form.Required(formFieldEmail, formFieldPassword, formFieldPassword2, formFieldUsername).
			MinLength(formFieldPassword, 8).
			MaxLength(formFieldPassword, 72).
			MinLength(formFieldPassword2, 8).
			MaxLength(formFieldPassword2, 72).
			MatchPass(formFieldPassword, formFieldPassword2).
			MinLength(formFieldUsername, 3)

		if !form.Valid() {
			app.render(c, "signup.html", &templateData{Form: form})
			return
		}

		email := c.Request.FormValue(formFieldEmail)
		username := c.Request.FormValue(formFieldUsername)
		password := c.Request.FormValue(formFieldPassword)

		// Check if email or username already used
		err = app.user.Validate(email, username)
		if err != nil {
			app.render(c, "signup.html", &templateData{Flash: err.Error()})
			return
		}

		// Add user to DB
		err = app.user.Add(email, password, username)
		if err != nil {
			app.render(c, "signup.html", &templateData{Flash: err.Error()})
			return
		}

		// 303 redirect (indicating POST to GET)
		http.Redirect(c.Writer, c.Request, "/login", http.StatusSeeOther)
		return
	}

	app.render(c, "signup.html", nil)
}

func (app *application) login(c *gin.Context) {

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
		form.Required(formFieldEmail, formFieldPassword)

		if !form.Valid() {
			app.errorLog.Printf("Validation failed: %v", form.Errors)
			app.render(c, "login.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		// check if user exists in DB
		user, err := app.user.GetUserByField(formFieldEmail, c.Request.FormValue(formFieldEmail))
		if err != nil {
			app.errorLog.Printf("Login failed: %v", err)
			app.render(c, "login.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		// check if pass from DB/input match
		if !checkPassword(user.Password, c.Request.FormValue(formFieldPassword)) {
			app.errorLog.Printf("Login failed due to incorrect password")
			app.render(c, "login.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		app.generateSession(c, user.ID)

		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
		return
	}

	app.render(c, "login.html", nil)

}

// Logoff deletes session and redirects to login page
func (app *application) logoff(c *gin.Context) {
	s, _ := app.store.Get(c.Request, app.sessionName)
	s.Options.MaxAge = -1

	fmt.Println("Successfully logged off")
	s.Save(c.Request, c.Writer)

	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
	// app.render("put in data for flash in template")
}

// TODO: Add functionality
func (app *application) myDecks(c *gin.Context) {
	app.render(c, "index.html", nil)
}
