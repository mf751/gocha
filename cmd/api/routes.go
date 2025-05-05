package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// router.HandlerFunc(http.MethodPost, "/auth", app.authenticate)
	router.HandlerFunc(http.MethodGet, "/ws", app.serveWS)
	return router
}
