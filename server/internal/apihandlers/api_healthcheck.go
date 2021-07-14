package apihandlers

import (
	"CloudScapes/pkg/logger"
	"CloudScapes/server/internal/dat"
	"context"
	"encoding/json"
	"net/http"
)

func HealthCheckGetHandler(db *dat.DB) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := db.PingDB(context.Background()); err != nil {
			logger.Log(logger.ERROR, "failed to ping", logger.Err(err))
			respond(rw, http.StatusInternalServerError, err.Error())
		}
		respond(rw, http.StatusOK, nil)
	}
}

func respond(w http.ResponseWriter, status int64, payload interface{}) bool {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log(logger.ERROR, "failed to marshal response", logger.Err(err))
		return false
	}

	if _, err = w.Write(marshalledPayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log(logger.ERROR, "failed to write response", logger.Err(err))
		return false
	}
	return true
}
