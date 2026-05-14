package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetFlash(c, n, v) gets session from store, then saves the session: session[n] = v
func (app *application) SetFlash(c *gin.Context, value string) {
	s, err := app.store.Get(c.Request, session_key)
	if err != nil {
		fmt.Println("Could not set flash due to: ", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("Adding value to flash: ", value)
	s.AddFlash(value)
	s.Save(c.Request, c.Writer)
}

func (app *application) GetFlash(c *gin.Context) string {
	fmt.Println("inside get flash")
	s, err := app.store.Get(c.Request, session_key)
	if err != nil {
		fmt.Println("Could not get flash due to: ", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	f := s.Flashes()
	fmt.Println("GetFlash flashes has: ", f)
	if f == nil {
		fmt.Println("flashes was empty")
		return ""
	}

	err = s.Save(c.Request, c.Writer)
	if err != nil {
		fmt.Println("Failed to save during getFlash: ", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}

	v, ok := f[0].(string)
	if !ok {
		fmt.Println("Failed to assert flash during getFlash: ", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	return v
}
