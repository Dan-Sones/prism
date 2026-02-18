package controller

import (
	"admin-service/internal/service"
	"net/http"
)

type EventsCatalogController struct {
	eventsCatalogService service.EventsCatalogServiceInterface
}

func NewEventsCatalogController(eventsCatalogService service.EventsCatalogServiceInterface) *EventsCatalogController {
	return &EventsCatalogController{
		eventsCatalogService: eventsCatalogService,
	}
}

func (e *EventsCatalogController) GetEventTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	searchQuery := r.URL.Query().Get("q")

	if searchQuery == "" {
		eventTypes, err := e.eventsCatalogService.GetEventTypes(ctx)
		if err != nil {
			WriteInternalServerError(w)
			return
		}
		WriteResponse(w, http.StatusOK, eventTypes)
		return
	}
	eventTypes, err := e.eventsCatalogService.SearchEventTypes(ctx, searchQuery)
	if err != nil {
		WriteInternalServerError(w)
		return
	}

	WriteResponse(w, http.StatusOK, eventTypes)
}
