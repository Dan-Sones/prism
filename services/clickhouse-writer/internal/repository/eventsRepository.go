package repository

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Dan-Sones/prismdbmodels/model"
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
		err := batch.Append(event.ExperimentDetails.ExperimentKey, event.ExperimentDetails.VariantKey, event.EventKey, event.UserDetails.ID, event.SentAt, event.ReceivedAt, event.GetStringProperties(), event.GetIntProperties(), event.GetFloatProperties())
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
