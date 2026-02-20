package controller

import (
	"admin-service/internal/problems"
	"admin-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		problems.NewBadRequestError("Request body is required").Write(w)
		return
	}

	var body model.EventType
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		problems.NewBadRequestError("Invalid request body").Write(w)
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

func (e *EventsCatalogController) IsFieldKeyAvailable(w http.ResponseWriter, r *http.Request) {
	eventTypeId := chi.URLParam(r, "eventTypeId")
	fieldKey := r.URL.Query().Get("fieldKey")

	if eventTypeId == "" || fieldKey == "" {
		problems.NewBadRequestError("eventTypeId and fieldKey are required").Write(w)
		return
	}

	if _, err := uuid.Parse(eventTypeId); err != nil {
		problems.NewBadRequestError("eventTypeId must be a valid UUID").Write(w)
		return
	}

	available, err := e.eventsCatalogService.IsFieldKeyAvailableForEventType(r.Context(), eventTypeId, fieldKey)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, map[string]bool{"available": available})
}

func (e *EventsCatalogController) DeleteEventType(w http.ResponseWriter, r *http.Request) {
	eventTypeId := chi.URLParam(r, "eventTypeId")

	if eventTypeId == "" {
		problems.NewBadRequestError("eventTypeId is required").Write(w)
		return
	}

	if _, err := uuid.Parse(eventTypeId); err != nil {
		problems.NewBadRequestError("eventTypeId must be a valid UUID").Write(w)
		return
	}

	err := e.eventsCatalogService.DeleteEventType(r.Context(), eventTypeId)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
