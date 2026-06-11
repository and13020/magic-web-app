package main

import (
	"fmt"
	"magic/repository"
	"magic/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	loggedInUserKey    = "user_id"
	formFieldEmail     = "email"
	formFieldPassword  = "password"
	formFieldUsername  = "username"
	formFieldPassword2 = "password2"
)

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

	err := c.Request.ParseForm()
	if err != nil {
		fmt.Println("Could not parse form! ", err)
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

func (app *application) signup(c *gin.Context) {

	if c.Request.Method == http.MethodPost {

		err := c.Request.ParseForm()
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		form := NewForm(c.Request.Form)

		// validations
		form.Required(formFieldEmail, formFieldPassword, formFieldPassword2, formFieldUsername).
			MinLength(formFieldPassword, 8).
			MaxLength(formFieldPassword, 255).
			MinLength(formFieldPassword2, 8).
			MaxLength(formFieldPassword2, 255).
			MatchPass(formFieldPassword, formFieldPassword2).
			MinLength(formFieldUsername, 3)

		if !form.Valid() {
			app.render(c, "signup.html", &templateData{Form: form}) // return form w/ error/s
			return
		}

		email := c.Request.FormValue(formFieldEmail)
		username := c.Request.FormValue(formFieldUsername)
		password := c.Request.FormValue(formFieldPassword)

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
		if !utils.CheckPassword(user.Password, c.Request.FormValue(formFieldPassword)) {
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
