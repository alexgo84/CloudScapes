package server

import (
	"CloudScapes/internal/server/dat"
	"context"
	"encoding/json"
	"net/http"
)

func healthCheckGetHandler(w http.ResponseWriter, r *http.Request) {
	if err := dat.PingDB(context.Background()); err != nil {
		respond(w, http.StatusInternalServerError, []byte(err.Error()))
	}
	respond(w, http.StatusOK, nil)

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
