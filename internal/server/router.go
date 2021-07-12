package server

import (
	"CloudScapes/internal/server/apihandlers"
	"CloudScapes/internal/server/dat"
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/logger"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func createRouter(db *dat.DB) *mux.Router {
	rootRouter := mux.NewRouter()
	rv1 := rootRouter.PathPrefix("/v1").Subrouter()

	// respond to not allowed same as not found to increase security
	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	rv1.NotFoundHandler = notFoundHandler
	rv1.MethodNotAllowedHandler = notFoundHandler

	// health check API (implement inline for simplicity since it doesnt use contextify)
	rv1.HandleFunc("/status/health",
		apihandlers.HealthCheckGetHandler(db)).
		Methods(http.MethodGet)

	// Accounts API
	rv1.HandleFunc("/accounts",
		contextify(db, authSession(apihandlers.AccountsGetHandler))).
		Methods(http.MethodGet)

	rv1.HandleFunc("/accounts",
		contextify(db, apihandlers.AccountsPostHandler)).
		Methods(http.MethodPost)

	return rootRouter
}

func contextify(db *dat.DB, h rqctx.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := rqctx.NewRequestContext(w, r)

		ctx.InitDBContext(db)

		defer func() {
			handlePanicRecovery(ctx)
		}()

		res := h(ctx)

		switch {
		case res.StatusCode >= 400:
			logger.Log(logger.INFO, "write ERROR response", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()), logger.Int64("status", int64(res.StatusCode)))
			if err := ctx.Rollback(); err != nil {
				logger.Log(logger.ERROR, "commit failed", zap.Error(err))
			}
			ctx.MarshalAndWrite([]byte(res.Err.Error()), res.StatusCode)

		default:
			logger.Log(logger.INFO, "write response", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()), logger.Int64("status", int64(res.StatusCode)))
			if err := ctx.Commit(); err != nil {
				logger.Log(logger.ERROR, "commit failed", zap.Error(err))
				return
			}
			ctx.MarshalAndWrite(res.Obj, res.StatusCode)
		}
	}
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
		logger.Log(logger.INFO, string(debug.Stack()))

		// rollback transaction
		rErr := ctx.Rollback()
		if err != nil {
			logger.Log(logger.ERROR, "failed to rollback error", zap.Error(rErr))
		}
	}
}
