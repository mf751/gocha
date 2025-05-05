package main

import "github.com/mf751/gocha/internal/jsonlog"

type config struct {
	port int
}

type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	app := &application{
		config: config{port: 5050},
	}
	app.serve()
}
