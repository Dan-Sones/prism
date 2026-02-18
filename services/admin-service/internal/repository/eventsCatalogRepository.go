package repository

import (
	"context"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsCatalogRepositoryInterface interface {
	CreateEventType(ctx context.Context, eventType model.EventType) error
	DeleteEventType(ctx context.Context, eventTypeId string) error
	GetEventTypes(ctx context.Context) ([]*model.EventType, error)
	SearchEventTypes(ctx context.Context, searchQuery string) ([]*model.EventType, error)
	IsFieldKeyAvailableForEventType(ctx context.Context, eventTypeId string, fieldKey string) (bool, error)
}

type EventsCatalogRepository struct {
	pgx *pgxpool.Pool
}

func NewEventsCatalogRepository(pgx *pgxpool.Pool) *EventsCatalogRepository {
	return &EventsCatalogRepository{
		pgx: pgx,
	}
}

func (e *EventsCatalogRepository) CreateEventType(ctx context.Context, eventType model.EventType) error {
	tx, err := e.pgx.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO prism.event_types (name, version, description) VALUES ($1, $2, $3) RETURNING id`

	var eventTypeId uuid.UUID
	err = tx.QueryRow(ctx, sql, eventType.Name, eventType.Version, eventType.Description).Scan(&eventTypeId)
	if err != nil {
		return err
	}

	for _, field := range eventType.Fields {
		sql = `INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type) VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(ctx, sql, eventTypeId, field.Name, field.FieldKey, field.DataType)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (e *EventsCatalogRepository) DeleteEventType(ctx context.Context, eventTypeId string) error {
	sql := `DELETE FROM prism.event_types WHERE id = $1`

	_, err := e.pgx.Exec(ctx, sql, eventTypeId)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventsCatalogRepository) GetEventTypes(ctx context.Context) ([]*model.EventType, error) {
	rows, err := e.pgx.Query(ctx, "SELECT id, name, version, description, created_at FROM prism.event_types")
	if err != nil {
		return nil, err
	}
	eventTypes, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.EventType])
	if err != nil {
		return nil, err
	}

	rows, err = e.pgx.Query(ctx, "SELECT id, event_type_id, name, field_key, data_type FROM prism.event_fields")
	if err != nil {
		return nil, err
	}
	fields, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.EventField])
	if err != nil {
		return nil, err
	}

	return combineEventTypesAndFields(eventTypes, fields), nil
}

func (e *EventsCatalogRepository) SearchEventTypes(ctx context.Context, searchQuery string) ([]*model.EventType, error) {
	rows, err := e.pgx.Query(ctx, "SELECT id, name, version, description, created_at FROM prism.event_types WHERE name ILIKE '%' || $1 || '%'", searchQuery)
	if err != nil {
		return nil, err
	}
	eventTypes, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.EventType])
	if err != nil {
		return nil, err
	}

	returnedIds := make([]uuid.UUID, len(eventTypes))
	for i, et := range eventTypes {
		returnedIds[i] = et.ID
	}

	rows, err = e.pgx.Query(ctx, "SELECT id, event_type_id, name, field_key, data_type FROM prism.event_fields WHERE event_type_id = ANY($1)", returnedIds)
	if err != nil {
		return nil, err
	}
	fields, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.EventField])
	if err != nil {
		return nil, err
	}

	return combineEventTypesAndFields(eventTypes, fields), nil
}

func combineEventTypesAndFields(eventTypes []*model.EventType, fields []model.EventField) []*model.EventType {
	fieldMap := make(map[uuid.UUID][]model.EventField)
	for _, f := range fields {
		fieldMap[f.EventTypeID] = append(fieldMap[f.EventTypeID], f)
	}
	for i := range eventTypes {
		eventTypes[i].Fields = fieldMap[eventTypes[i].ID]
	}

	return eventTypes
}

func (e *EventsCatalogRepository) IsFieldKeyAvailableForEventType(ctx context.Context, eventTypeId string, fieldKey string) (bool, error) {
	//TODO implement me
	panic("implement me")
	// this will be used when creating a new field in an event type, the frontend wil query before allowing the form to be submitted.
}
