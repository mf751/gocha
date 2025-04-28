package main

import (
	"net/http"
)

func (app *application) serveWS(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Running!"))
}
