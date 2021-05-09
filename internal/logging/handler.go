package logging

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("incoming request: ", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
