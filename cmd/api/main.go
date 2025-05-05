package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/mf751/gocha/internal/data"
	"github.com/mf751/gocha/internal/jsonlog"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Modles
}

func main() {
	var cfg config
	cfg.port = 5050
	cfg.db.dsn = ""
	cfg.db.maxIdleConns = 25
	cfg.db.maxOpenConns = 25
	cfg.db.maxIdleTime = "15m"
	cfg.cors.trustedOrigins = []string{"http://localhost:5173"}
	cfg.limiter.rps = 5
	cfg.limiter.burst = 8
	cfg.limiter.enabled = true

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("Database connection established", nil)

	app := &application{
		config: cfg,
		models: data.NewModels(db),
		logger: logger,
	}

	app.serve()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	return db, err
}
