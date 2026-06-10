package main

import (
	"log"
	"net/http"
	"time"
)

func (app *application) Serve() {
	s := http.Server{
		Addr:         ":8080",          // config
		ReadTimeout:  5 * time.Second,  // config
		WriteTimeout: 10 * time.Second, // config
		Handler:      app.routes(),
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
