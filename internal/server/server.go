package server

import (
	"CloudScapes/internal/server/dat"
	"CloudScapes/internal/server/rqctx"
	l "CloudScapes/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
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

func handlePanicRecovery(ctx *rqctx.Context) {
	if err := recover(); err != nil {

		// report error
		e, ok := err.(error)
		if !ok {
			e = fmt.Errorf("%#v\n%s", err, string(debug.Stack()))
		}
		ctx.ReportError("Recovered handler panic", zap.Error(e))

		// print stack
		l.Log(l.INFO, string(debug.Stack()))

		// rollback transaction
		rErr := ctx.Rollback()
		if err != nil {
			l.Log(l.ERROR, "failed to rollback error", zap.Error(rErr))
		}
	}
}
