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
		s, err := app.store.Get(c.Request, session_key)
		if err != nil {
			fmt.Println("Could not decode session: ", err)
			return
		}
		if _, ok := s.Values["userID"]; !ok {
			fmt.Println("session value USERID not found: ", ok)
			fmt.Println("session contents: ", s)
			http.Redirect(c.Writer, c.Request, "/login", http.StatusFound)
			return
		}

		cookie, err := c.Request.Cookie(session_key)
		if err == http.ErrNoCookie {
			fmt.Println("Cookie not present: ", err)
			http.Redirect(c.Writer, c.Request, "/login", http.StatusUnauthorized)
			return
		} else if err != nil {
			fmt.Println("Could not get cookie from request: ", err)
			http.Redirect(c.Writer, c.Request, "/login", http.StatusInternalServerError)
			return
		}

		if u, err := app.user.GetUserByField("id", cookie.Value); u.ID != cookie.Value || err != nil {
			fmt.Println("Cookie validation failed: ", err)
			http.Redirect(c.Writer, c.Request, "/login", http.StatusUnauthorized)
			return
		}
		fmt.Println("User is logged in, session exists")

		c.AddParam("auth", "true")
		fmt.Println("in middleware : ", c.Value("auth"))
		c.Next()
	}

}
