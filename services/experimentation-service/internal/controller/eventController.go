package controller

import (
	model2 "experimentation-service/internal/model/graph"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type EventController struct {
	eventService *service.EventService
}

func NewEventController(eventService *service.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

func (e *EventController) GetEventUsageOverPeriod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventKey := chi.URLParam(r, "eventKey")
	graphScaleParam := r.URL.Query().Get("graphScale")

	graphScale, ok := model2.ParseGraphTimeScale(graphScaleParam)
	if !ok {
		problems.NewBadRequestError("Invalid graphScale value, must be one of: " + model2.GetListOfGraphTimeScales()).Write(w)
		return
	}

	if eventKey == "" {
		problems.NewBadRequestError("eventTypeId is required").Write(w)
		return
	}

	dataPoints, err := e.eventService.GetEventUsageOverPeriod(ctx, graphScale, eventKey)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, dataPoints)
}
