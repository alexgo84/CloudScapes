package apihandlers

import (
	"CloudScapes/server/internal/dat"
	"context"
	"encoding/json"
	"net/http"
)

func HealthCheckGetHandler(db *dat.DB) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := db.PingDB(context.Background()); err != nil {
			respond(rw, http.StatusInternalServerError, []byte(err.Error()))
		}
		respond(rw, http.StatusOK, nil)
	}
}

func respond(w http.ResponseWriter, status int64, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(marshalledPayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
