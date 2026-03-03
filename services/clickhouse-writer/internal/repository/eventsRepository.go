package repository

import (
	"clickhouse-writer/internal/model"
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type EventsRepository interface {
	InsertBatch(ctx context.Context, events []model.DownstreamEvent) error
}

type EventsRepositoryClickhouse struct {
	connection driver.Conn
}

func NewEventsRepositoryClickhouse(connection driver.Conn) *EventsRepositoryClickhouse {
	return &EventsRepositoryClickhouse{
		connection: connection,
	}
}

func (r *EventsRepositoryClickhouse) InsertBatch(ctx context.Context, events []model.DownstreamEvent) error {
	insertCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	batch, err := r.connection.PrepareBatch(insertCtx, "INSERT INTO events")
	if err != nil {
		return err
	}

	for _, event := range events {
		err := batch.Append(event.EventKey, event.UserDetails.ID, event.SentAt, event.ReceivedAt, getStringProperties(event), getIntProperties(event), getFloatProperties(event))
		if err != nil {
			return err
		}
	}

	return batch.Send()
}

func getStringProperties(event model.DownstreamEvent) map[string]string {
	stringProperties := make(map[string]string)
	for key, value := range event.Properties {
		if value.DataType == model.OutboundEventFieldDataTypeString {
			stringProperties[key] = value.Value.(string)
		}
	}
	return stringProperties
}

func getFloatProperties(event model.DownstreamEvent) map[string]float64 {
	floatProperties := make(map[string]float64)
	for key, value := range event.Properties {
		if value.DataType == model.OutboundEventFieldDataTypeFloat {
			floatProperties[key] = value.Value.(float64)
		}
	}
	return floatProperties
}

func getIntProperties(event model.DownstreamEvent) map[string]int64 {
	intProperties := make(map[string]int64)
	for key, value := range event.Properties {
		if value.DataType == model.OutboundEventFieldDataTypeInt {
			intProperties[key] = value.Value.(int64)
		}
	}
	return intProperties
}
