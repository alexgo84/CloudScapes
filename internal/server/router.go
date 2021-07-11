package server

import (
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/logger"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func createRouter() *mux.Router {
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

	// Accounts API
	rv1.HandleFunc("/accounts",
		contextify(authSession(accountsGetHandler))).
		Methods(http.MethodGet)

	rv1.HandleFunc("/accounts",
		contextify(accountsPostHandler)).
		Methods(http.MethodPost)

	return rootRouter
}

func contextify(h rqctx.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := rqctx.NewRequestContext(w, r)

		ctx.InitDBTransaction()

		defer func() {
			handlePanicRecovery(ctx)
		}()

		res := h(ctx)

		switch {
		case res.StatusCode == http.StatusNoContent:
			ctx.MarshalAndWrite([]byte{}, res.StatusCode)
		case res.StatusCode >= 400:
			ctx.MarshalAndWrite([]byte(res.Err.Error()), res.StatusCode)
		default:
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
