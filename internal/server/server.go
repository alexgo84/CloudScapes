package server

import (
	l "CloudScapes/pkg/logger"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	l.Log(l.INFO, "request received on server", l.Str("method", r.Method), l.Str("URL", r.URL.String()))
}

func Run() error {
	http.HandleFunc("/", handler)

	l.Log(l.INFO, "listening on port 8080")
	return http.ListenAndServe(":8080", nil)
}
