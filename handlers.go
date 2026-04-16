package main

import (
	"fmt"
	"magic/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// mux.HandleFunc("/", app.home)
// mux.HandleFunc("/search/{name}", app.card.search) // {name} is a path variable that can be accessed in the handler function
// mux.HandleFunc("/random", app.card.random)
// mux.HandleFunc("/update", app.card.updateCard)
// mux.HandleFunc("/delete", app.card.deleteCard)

const loggedInUserKey = "user_id"

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
