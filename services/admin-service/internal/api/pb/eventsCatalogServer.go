package pb

import (
	"admin-service/internal/grpc/generated/events_catalog/v1"
	"admin-service/internal/service"
	"context"

	"github.com/Dan-Sones/prismdbmodels/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventsCatalogServer struct {
	events_catalog.UnimplementedEventsCatalogServiceServer
	eventsCatalogService service.EventsCatalogServiceInterface
}

func NewEventsCatalogServer(eventsCatalogService service.EventsCatalogServiceInterface) *EventsCatalogServer {
	return &EventsCatalogServer{
		eventsCatalogService: eventsCatalogService,
	}
}

func (e EventsCatalogServer) GetEventTypeByKey(ctx context.Context, in *events_catalog.GetEventTypeByKeyRequest) (*events_catalog.EventType, error) {
	eventType, err := e.eventsCatalogService.GetEventTypeByKey(ctx, in.EventKey)
	if err != nil {
		return nil, err
	}

	return &events_catalog.EventType{
		Id:          eventType.ID.String(),
		Name:        eventType.Name,
		EventKey:    eventType.EventKey,
		Version:     int32(eventType.Version),
		Description: eventType.Description,
		CreatedAt:   timestamppb.New(eventType.CreatedAt),
		Fields:      convertEventFields(eventType.Fields),
	}, nil
}

func convertEventFields(fields []model.EventField) []*events_catalog.EventField {
	converted := make([]*events_catalog.EventField, len(fields))
	for i, field := range fields {
		converted[i] = &events_catalog.EventField{
			Id:       field.ID.String(),
			Name:     field.Name,
			FieldKey: field.FieldKey,
			DataType: convertDataType(field.DataType),
		}
	}
	return converted
}

func convertDataType(dataType model.DataType) events_catalog.DataType {
	switch dataType {
	case model.DataTypeString:
		return events_catalog.DataType_DATA_TYPE_STRING
	case model.DataTypeInt:
		return events_catalog.DataType_DATA_TYPE_INT
	case model.DataTypeFloat:
		return events_catalog.DataType_DATA_TYPE_FLOAT
	case model.DataTypeBoolean:
		return events_catalog.DataType_DATA_TYPE_BOOL
	case model.DataTypeTimestamp:
		return events_catalog.DataType_DATA_TYPE_TIMESTAMP
	default:
		return events_catalog.DataType_DATA_TYPE_UNSPECIFIED
	}
}
