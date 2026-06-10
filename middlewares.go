package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: need unique key for each user to validate auth
// sessionMiddleware reads session id from cookie
// Then we verify the session id (db/JWT secret)
// Store auth data in context to avoid re-validating in each handler
func (app *application) sessionMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// TODO: If user not logged in, skip and go to next handler
		s, err := app.store.Get(c.Request, app.sessionName) // session key
		if err != nil {
			fmt.Println("Could not decode session: ", err)
			return
		}
		userID, ok := s.Values[loggedInUserKey].(string)
		if !ok || userID == "" {
			fmt.Println("session value USERID not found or invalid")
			fmt.Println("session contents: ", s)
			http.Redirect(c.Writer, c.Request, "/login", http.StatusFound)
			return
		}

		if u, err := app.user.GetUserByField("id", userID); err != nil || u.ID != userID {
			fmt.Println("User validation failed: ", err)
			http.Redirect(c.Writer, c.Request, "/login", http.StatusUnauthorized)
			return
		}
		fmt.Println("User is logged in, session exists")

		c.AddParam("auth", "true")
		fmt.Println("in middleware : ", c.Value("auth"))
		c.Next()
	}

}
