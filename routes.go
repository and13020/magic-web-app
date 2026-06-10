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
	public.Match(GetAndPost, "/signup", app.signup)
	public.Match(GetAndPost, "/login", app.login)
	public.GET("/logoff", app.logoff)
	public.GET("/random", app.random)

	private := r.Group("/", app.sessionMiddleware())
	private.GET("/mydecks", app.myDecks)

	return r
}
