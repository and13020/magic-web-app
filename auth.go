package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// isAuthenticated accepts a *http.Request.
// It checks store for existing session. Returns if authenticated
func (app *application) isAuthenticated(r *http.Request) bool {
	s, err := app.store.Get(r, app.sessionName)
	if err != nil {
		fmt.Println("Could not access store: ", err)
		return false
	}

	userID, ok := s.Values[loggedInUserKey].(string)
	if !ok || userID == "" {
		return false
	}

	_, err = app.user.GetUserByField("id", userID)
	if err != nil {
		fmt.Printf("Could not verify user id from session: %v\n", err)
		return false
	}

	return true
}

func (app *application) generateSession(c *gin.Context, uID string) {
	// Create session and store user id
	session, err := app.store.Get(c.Request, app.sessionName)
	if err != nil {
		fmt.Println("Failed to store.Get to create new session: ", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values[loggedInUserKey] = uID

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		fmt.Println("session failed to save: ", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Session created: ", session)
	app.SetFlash(c, "Successfully logged in")
}
