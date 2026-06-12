package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func failAuth(c *gin.Context, s *sessions.Session) {
	s.Options.MaxAge = -1
	_ = s.Save(c.Request, c.Writer)
	http.Redirect(c.Writer, c.Request, "/login", http.StatusUnauthorized)
}

func (app *application) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate session from cookiestore
		s, err := app.store.Get(c.Request, app.sessionName)
		if err != nil {
			failAuth(c, s)
			return
		}

		// Validate user ID from db
		uID, ok := s.Values[UserIdKey].(string)
		if !ok || uID == "" {
			app.errorLog.Printf("Auth failed to retrieve session: %v", ok)
			failAuth(c, s)
			return
		}
		_, err = app.user.GetUserByField(UserIdKey, uID)
		if err != nil {
			app.errorLog.Printf("Auth failed to retrieve user due to %v", err)
			failAuth(c, s)
			return
		}

		c.Next()
	}
}
