package main

import (
	"context"
	"database/sql"
	"flag"
	"gopherinha/internal/data"
	"gopherinha/internal/jsonlog"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
    port int
    env string
    db   struct {
        dsn          string
        maxOpenConns int
        maxIdleConns int
        maxIdleTime  string
    }
}

type application struct {
    config config 
    logger *jsonlog.Logger
    models data.Models 
    wg sync.WaitGroup
}

func main() {
    var cfg config 

    flag.IntVar(&cfg.port, "port", 4000, "API server port")
    flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
    flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

    flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
    flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
    flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

    flag.Parse()

    logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

    db, err := openDB(cfg)
    if err != nil {
        logger.PrintFatal(err, nil)
    }
    defer db.Close()

    app := &application{
        config: cfg,
        models: data.NewModels(db),
        logger: logger,
    }

    err = app.serve()
    if err != nil {
        logger.PrintFatal(err, nil)
    }

    logger.PrintFatal(err, nil)
}


func openDB(cfg config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.db.dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(cfg.db.maxOpenConns)
    db.SetMaxIdleConns(cfg.db.maxIdleConns)

    duration, err := time.ParseDuration(cfg.db.maxIdleTime)
    if err != nil {
        return nil, err
    }

    db.SetConnMaxIdleTime(duration)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = db.PingContext(ctx)
    if err != nil {
        return nil, err
    }

    return db, nil
}
