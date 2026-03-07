package http

import (
	"experimentation-service/internal/controller"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controllers struct {
	ExperimentController    *controller.ExperimentController
	EventsCatalogController *controller.EventsCatalogController
	EventController         *controller.EventController
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
			r.Get("/event-key-available", c.EventsCatalogController.IsEventKeyAvailable)

			r.Route("/byKey/{eventKey}", func(r chi.Router) {
				r.Get("/usage", c.EventController.GetEventUsageOverPeriod)
			})

			r.Route("/{eventTypeId}", func(r chi.Router) {
				r.Get("/", c.EventsCatalogController.GetEventType)
				r.Delete("/", c.EventsCatalogController.DeleteEventType)
				r.Get("/field-key-available", c.EventsCatalogController.IsFieldKeyAvailable)
			})
		})
	})
}
