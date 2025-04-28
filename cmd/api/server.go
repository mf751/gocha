package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) serve() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("starting server on port: %v", app.config.port)

	log.Fatal(srv.ListenAndServe())
}
