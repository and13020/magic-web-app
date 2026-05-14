package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {

	GetAndPost := []string{"POST", "GET"}
	r := gin.Default() // default includes logging and recovery middleware
	// gin.SetMode(gin.ReleaseMode)

	r.Static(app.publicPath, app.publicPath)

	public := r.Group("/")
	public.GET("/", app.home)
	public.GET("/search", app.getCards)
	public.POST("/search", app.getCardsForm)
	public.Match(GetAndPost, "/signup", app.Signup)
	public.Match(GetAndPost, "/login", app.Login)
	public.GET("/logoff", app.Logoff)
	public.GET("/random", app.random)

	// private := r.Group("/", app.sessionMiddleware())
	// removing this group/middleware as we can just validate for auth on each render

	// proper flow for JWT:
	// 1. user creates account
	// 2. user signs in (we validate credentials)
	// 3. new JWT token created w/ users identifiers and expiration timestamp
	// 4. encode header/payload sign them w/ secret key to create signature (JWT created)
	// 5. server returns JWT as response
	// 6. client returns JWT on subsequent requsts
	// 7. server validates token by creating new signature, comparing w/ existing signature (and checking expiration too)

	// Routes

	// problem:
	/*	we need to enable session data for cache
		we also need to wrap our handlers with the session middleware to access session data in our handlers
		we can wrap our handlers with the session middleware using gin.WrapH and gin.WrapF
	*/

	// r.GET("/random", app.abcd)
	// r.PUT("/update", app.card.UpdateCard)
	// r.DELETE("/delete", app.card.DeleteCard)

	return r
}
