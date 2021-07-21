package server

import (
	"CloudScapes/pkg/shared/redis"
	"CloudScapes/server/internal/apihandlers"
	"CloudScapes/server/internal/dat"
	"net/http"

	"github.com/gorilla/mux"
)

func createRouter(db *dat.DB, ps *redis.PubSubClient) *mux.Router {
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
		contextify(db, ps, authSession(apihandlers.AccountsGetHandler))).
		Methods(http.MethodGet)

	rv1.HandleFunc("/accounts",
		apihandlers.AccountsPostHandler(db)).
		Methods(http.MethodPost)

	// Users API
	rv1.HandleFunc("/users",
		contextify(db, ps, apihandlers.UsersPostHandler)).
		Methods(http.MethodPost)

	rv1.HandleFunc("/users",
		contextify(db, ps, apihandlers.UsersGetHandler)).
		Methods(http.MethodGet)

	// Clusters API
	rv1.HandleFunc("/clusters",
		contextify(db, ps, apihandlers.ClustersPostHandler)).
		Methods(http.MethodPost)

	rv1.HandleFunc("/clusters",
		contextify(db, ps, apihandlers.ClustersGetHandler)).
		Methods(http.MethodGet)

	rv1.HandleFunc("/clusters/{clusterId}",
		contextify(db, ps, apihandlers.ClustersDeleteHandler)).
		Methods(http.MethodDelete)

	// Plans API
	rv1.HandleFunc("/plans",
		contextify(db, ps, apihandlers.PlansPostHandler)).
		Methods(http.MethodPost)

	rv1.HandleFunc("/plans",
		contextify(db, ps, apihandlers.PlansGetHandler)).
		Methods(http.MethodGet)

	rv1.HandleFunc("/plans/{planId}",
		contextify(db, ps, apihandlers.PlansPutHandler)).
		Methods(http.MethodPut)

	rv1.HandleFunc("/plans/{planId}",
		contextify(db, ps, apihandlers.PlansDeleteHandler)).
		Methods(http.MethodDelete)

	// Deployments API
	rv1.HandleFunc("/deployments",
		contextify(db, ps, apihandlers.DeploymentsPostHandler)).
		Methods(http.MethodPost)

	rv1.HandleFunc("/deployments",
		contextify(db, ps, apihandlers.DeploymentsGetHandler)).
		Methods(http.MethodGet)

	rv1.HandleFunc("/deployments/{deploymentId}",
		contextify(db, ps, apihandlers.DeploymentsPutHandler)).
		Methods(http.MethodPut)

	rv1.HandleFunc("/deployments/{deploymentId}",
		contextify(db, ps, apihandlers.DeploymentsDeleteHandler)).
		Methods(http.MethodDelete)

	return rootRouter
}
