package server

import (
	"CloudScapes/internal/server/dat"
	l "CloudScapes/pkg/logger"
	"context"
	"net/http"
	"time"
)

func Run() error {
	l.Log(l.INFO, "initializing database")

	db, err := dat.NewDB(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Log(l.ERROR, "database close failed", l.Err(err))
		}
	}()

	l.Log(l.INFO, "running all database migrations")
	if err := db.RunMigrations(context.Background()); err != nil {
		return err
	}

	s := &http.Server{
		Handler:      createRouter(db),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	l.Log(l.INFO, "serving requests on port 8080")
	return s.ListenAndServe()
}
