package server

import (
	"CloudScapes/internal/server/dat"
	l "CloudScapes/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	rootRouter := mux.NewRouter()
	rv1 := rootRouter.PathPrefix("/v1").Subrouter()

	// respond to not allowed same as not found to increase security
	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	rv1.NotFoundHandler = notFoundHandler
	rv1.MethodNotAllowedHandler = notFoundHandler

	// health check API
	rv1.HandleFunc("/status/health",
		healthCheckGetHandler).
		Methods(http.MethodGet)

	return &http.Server{
		Handler:      rootRouter,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
