package main

type config struct {
	port int
}

type application struct {
	config config
}

func main() {
	app := &application{
		config: config{port: 5050},
	}
	app.serve()
}
