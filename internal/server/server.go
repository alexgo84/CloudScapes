package server

import (
	"CloudScapes/internal/server/dat"
	l "CloudScapes/pkg/logger"
	"context"
	"net/http"
	"time"
)

func Run() error {
	if err := dat.InitDB(context.Background()); err != nil {
		return err
	}
	defer func() {
		if err := dat.CloseDB(context.Background()); err != nil {
			l.Log(l.ERROR, "database close failed", l.Err(err))
		}
	}()

	s := createServer()
	l.Log(l.INFO, "listening on port 8080")
	return s.ListenAndServe()
}

func createServer() *http.Server {

	return &http.Server{
		Handler:      createRouter(),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
