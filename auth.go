package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values[loggedInUserKey] = uID

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	app.SetFlash(c, "Successfully logged in")
}

// checkPassword accepts a hashed password and plaintext password, compares if they match.
// Returns true on a match.
func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	return false
}
