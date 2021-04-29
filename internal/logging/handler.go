package logging

import (
	"log"
	"net/http"
)

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming request: ", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
