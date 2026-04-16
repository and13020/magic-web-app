package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {

	// Set up router
	r := gin.Default() // default includes logging and recovery middleware
	// gin.SetMode(gin.ReleaseMode)

	r.Static(app.publicPath, app.publicPath)

	public := r.Group("/v1")
	public.GET("/", app.home)
	public.GET("/random", app.random)
	public.GET("/search", app.getCards)
	public.POST("/search", app.getCardsForm)

	// proper flow for JWT:
	// 1. user creates account
	// 2. user signs in (we validate credentials)
	// 3. new JWT token created w/ users identifiers and expiration timestamp
	// 4. encode header/payload sign them w/ secret key to create signature (JWT created)
	// 5. server returns JWT as response
	// 6. client returns JWT on subsequent requsts
	// 7. server validates token by creating new signature, comparing w/ existing signature (and checking expiration too)
	private := r.Group("/v2")
	private.GET("/", app.home)

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
