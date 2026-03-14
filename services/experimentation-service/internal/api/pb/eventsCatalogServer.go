package pb

import (
	"context"
	"experimentation-service/internal/grpc/generated/events_catalog/v1"
	"experimentation-service/internal/service"

	"github.com/Dan-Sones/prismdbmodels/model/event"
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

func (e EventsCatalogServer) GetEventTypeByKey(ctx context.Context, in *events_catalog.GetEventTypeByKeyRequest) (*events_catalog.GetEventTypeByKeyResponse, error) {
	eventType, err := e.eventsCatalogService.GetEventTypeByKey(ctx, in.EventKey)
	if err != nil {
		return nil, err
	}

	return &events_catalog.GetEventTypeByKeyResponse{
		EventType: &events_catalog.EventType{
			Id:          eventType.ID.String(),
			Name:        eventType.Name,
			EventKey:    eventType.EventKey,
			Version:     int32(eventType.Version),
			Description: eventType.Description,
			CreatedAt:   timestamppb.New(eventType.CreatedAt),
			Fields:      convertEventFields(eventType.Fields),
		},
	}, nil
}

func convertEventFields(fields []event.EventField) []*events_catalog.EventField {
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

func convertDataType(dataType event.DataType) events_catalog.DataType {
	switch dataType {
	case event.DataTypeString:
		return events_catalog.DataType_DATA_TYPE_STRING
	case event.DataTypeInt:
		return events_catalog.DataType_DATA_TYPE_INT
	case event.DataTypeFloat:
		return events_catalog.DataType_DATA_TYPE_FLOAT
	case event.DataTypeBoolean:
		return events_catalog.DataType_DATA_TYPE_BOOL
	case event.DataTypeTimestamp:
		return events_catalog.DataType_DATA_TYPE_TIMESTAMP
	default:
		return events_catalog.DataType_DATA_TYPE_UNSPECIFIED
	}
}
