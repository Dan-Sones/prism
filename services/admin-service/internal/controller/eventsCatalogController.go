package controller

import (
	"admin-service/internal/problems"
	"admin-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type EventsCatalogController struct {
	eventsCatalogService service.EventsCatalogServiceInterface
}

func NewEventsCatalogController(eventsCatalogService service.EventsCatalogServiceInterface) *EventsCatalogController {
	return &EventsCatalogController{
		eventsCatalogService: eventsCatalogService,
	}
}

func (e *EventsCatalogController) CreateEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil {
		err := problems.NewBadRequestError()
		err.Write(w)
		return
	}

	var body model.EventType
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		problems.NewBadRequestError().Write(w)
		return
	}

	err, violations := e.eventsCatalogService.CreateEventType(ctx, body)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	if violations != nil && len(violations) > 0 {
		problems.NewValidationError(violations).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *EventsCatalogController) GetEventTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	searchQuery := r.URL.Query().Get("q")

	if searchQuery == "" {
		eventTypes, err := e.eventsCatalogService.GetEventTypes(ctx)
		if err != nil {
			problems.NewInternalServerError().Write(w)
			return
		}
		WriteResponse(w, http.StatusOK, eventTypes)
		return
	}
	eventTypes, err := e.eventsCatalogService.SearchEventTypes(ctx, searchQuery)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, eventTypes)
}
