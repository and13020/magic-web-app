package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) sessionMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// If user not logged in, skip and go to next handler
		if !app.session.Exists(c.Request, loggedInUserKey) {
			c.Next()
			return
		}

		// if session exists, check if user is logged in
		// If user is logged in - set some flag true
		fmt.Println("User is logged in, session exists")

		// Get the session for the current request
		app.session.Enable(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		c.Next()
	}

}
