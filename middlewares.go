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

		// TODO: if we update keys or smth we should remove invalid cookies
		// when users access any private routes requiring this sessionMiddle()
		// which then reroutes to login anew

		// problem is the MaxAge field doesn't cause client to remove the cookie or do anything..
		// so we'd have to read the time and compare to present time or smth?

		userID, ok := s.Values[loggedInUserKey].(string)
		if !ok || userID == "" {
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
