package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controllers struct {
}

func RegisterRoutes(router *chi.Mux, c Controllers) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		})
	})
}
