package server

import (
	"CloudScapes/internal/server/apihandlers"
	"CloudScapes/internal/server/dat"
	"net/http"

	"github.com/gorilla/mux"
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

	// Users API
	rv1.HandleFunc("/users",
		contextify(db, apihandlers.UsersPostHandler)).
		Methods(http.MethodPost)

	rv1.HandleFunc("/users",
		contextify(db, apihandlers.UsersGetHandler)).
		Methods(http.MethodGet)

	return rootRouter
}
