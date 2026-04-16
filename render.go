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
	app.tp.Render(c.Writer, filename, app.defaultTemplateData(data, c.Request))
}

func (app *application) defaultTemplateData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}
	// data.Flash = "flash was here"	// flash works
	data.IsAuthenticated = app.isAuthenticated(r) // add auth later

	return data
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuth, exists := r.Context().Value("user_id").(bool)

	if !exists {
		return false
	}

	return isAuth
}
