package http

import (
	"admin-service/internal/controller"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controllers struct {
	ExperimentController    *controller.ExperimentController
	EventsCatalogController *controller.EventsCatalogController
}

func RegisterRoutes(router *chi.Mux, c Controllers) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		})

		r.Post("/experiments", c.ExperimentController.CreateExperiment)

		r.Route("/events-catalog", func(r chi.Router) {
			r.Get("/", c.EventsCatalogController.GetEventTypes)
			r.Post("/", c.EventsCatalogController.CreateEventType)
		})
	})
}
