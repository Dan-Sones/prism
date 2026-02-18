package service

import (
	"admin-service/internal/repository"
	"context"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type EventsCatalogServiceInterface interface {
	CreateEventType(ctx context.Context, eventType model.EventType) error
	DeleteEventType(ctx context.Context, eventTypeId string) error
	GetEventTypes(ctx context.Context) ([]*model.EventType, error)
	SearchEventTypes(ctx context.Context, searchQuery string) ([]*model.EventType, error)
}

type EventsCatalogService struct {
	eventsCatalogRepository repository.EventsCatalogRepositoryInterface
	logger                  *slog.Logger
}

func NewEventsCatalogService(eventsCatalogRepository repository.EventsCatalogRepositoryInterface, logger *slog.Logger) *EventsCatalogService {
	return &EventsCatalogService{
		eventsCatalogRepository: eventsCatalogRepository,
		logger:                  logger,
	}
}

func (e *EventsCatalogService) CreateEventType(ctx context.Context, eventType model.EventType) error {
	panic("implement me")
}

func (e *EventsCatalogService) DeleteEventType(ctx context.Context, eventTypeId string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventsCatalogService) GetEventTypes(ctx context.Context) ([]*model.EventType, error) {
	eventTypes, err := e.eventsCatalogRepository.GetEventTypes(ctx)
	if err != nil {
		e.logger.Error("Error fetching event types", "error", err)
		return nil, err
	}
	return eventTypes, nil
}

func (e *EventsCatalogService) SearchEventTypes(ctx context.Context, searchQuery string) ([]*model.EventType, error) {
	eventTypes, err := e.eventsCatalogRepository.SearchEventTypes(ctx, searchQuery)
	if err != nil {
		e.logger.Error("Error searching event types", "error", err, "searchQuery", searchQuery)
		return nil, err
	}
	return eventTypes, nil
}
