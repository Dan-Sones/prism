package service

import (
	"admin-service/internal/problems"
	"admin-service/internal/repository"
	"admin-service/internal/validators"
	"context"
	"errors"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type EventsCatalogServiceInterface interface {
	CreateEventType(ctx context.Context, eventType model.EventType) (error, []problems.Violation)
	DeleteEventType(ctx context.Context, eventTypeId string) error
	GetEventTypes(ctx context.Context) ([]*model.EventType, error)
	SearchEventTypes(ctx context.Context, searchQuery string) ([]*model.EventType, error)
	IsFieldKeyAvailableForEventType(ctx context.Context, eventTypeId string, fieldKey string) (bool, error)
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

func (e *EventsCatalogService) CreateEventType(ctx context.Context, eventType model.EventType) (error, []problems.Violation) {
	violations := validators.ValidateEventType(eventType)
	if len(violations) > 0 {
		return nil, violations
	}

	err := e.eventsCatalogRepository.CreateEventType(ctx, eventType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var violation problems.Violation
			switch pgErr.ConstraintName {
			case "unique_event_type_name":
				violation = problems.Violation{Field: "name", Message: "An event type with this name already exists"}
			case "unique_event_type_field_key":
				violation = problems.Violation{Field: "fieldKey", Message: "A field with this key already exists for this event type"}
			default:
				violation = problems.Violation{Field: "unknown", Message: "A unique constraint violation occurred"}
			}
			e.logger.Error(violation.Message, "eventType", eventType)
			return nil, []problems.Violation{violation}
		}

		e.logger.Error("Error creating event type", "error", err, "eventType", eventType)
		return err, nil
	}

	return nil, nil
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

func (e *EventsCatalogService) IsFieldKeyAvailableForEventType(ctx context.Context, eventTypeId string, fieldKey string) (bool, error) {
	available, err := e.eventsCatalogRepository.IsFieldKeyAvailableForEventType(ctx, eventTypeId, fieldKey)
	if err != nil {
		e.logger.Error("Error checking field key availability", "error", err, "eventTypeId", eventTypeId, "fieldKey", fieldKey)
		// default to it not being available if there's an error to be safe.
		return false, err
	}

	return available, nil
}
