package logging

import (
	"net/http"
)

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Info.Println("incoming request: ", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
