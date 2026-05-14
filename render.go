package main

import (
	"fmt"
	"magic/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type templateData struct {
	Form            *Form
	IsAuthenticated bool
	Flash           string
	Cards           []repository.Card
	Card            repository.Card
	NextLink        string
	PrevLink        string
}

func (app *application) render(c *gin.Context, filename string, data *templateData) {
	if app.tp == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println("Template renderer not initialized")
	}
	app.tp.Render(c.Writer, filename, app.defaultTemplateData(c, data))
}

func (app *application) defaultTemplateData(c *gin.Context, data *templateData) *templateData {
	if data == nil {
		data = &templateData{}
	}
	// data.Flash = "flash was here"                 // flash works
	if f := app.GetFlash(c); f != "" {
		data.Flash = f
	}
	data.IsAuthenticated = app.isAuthenticated(c.Request) // add auth later

	return data
}

// isAuthenticated accepts a *http.Request.
// It checks store for existing session. Returns if authenticated
func (app *application) isAuthenticated(r *http.Request) bool {

	s, err := app.store.Get(r, session_key)
	if err != nil {
		fmt.Println("Could not access store: ", err)
		return false
	}
	auth, ok := s.Values[loggedInUserKey].(bool)
	if !ok || !auth {
		fmt.Printf("Could not verify auth: %t", auth)
		return false
	}

	return auth
}
