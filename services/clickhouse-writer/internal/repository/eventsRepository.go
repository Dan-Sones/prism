package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Dan-Sones/prismdbmodels/model"
)

type EventsRepository interface {
	InsertBatch(ctx context.Context, events []model.DownstreamEvent) error
}

type EventsRepositoryClickhouse struct {
	connection driver.Conn
	logger     *slog.Logger
}

func NewEventsRepositoryClickhouse(connection driver.Conn, logger *slog.Logger) *EventsRepositoryClickhouse {
	return &EventsRepositoryClickhouse{
		connection: connection,
		logger:     logger,
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
		stringProps, stringErr := event.GetStringProperties()
		intProps, intErr := event.GetIntProperties()
		floatProps, floatErr := event.GetFloatProperties()
		// Drop the mistyped properties but keep the event: failing here would block the
		// whole microbatch from committing and redeliver it forever.
		if err := errors.Join(stringErr, intErr, floatErr); err != nil {
			r.logger.Warn("dropping event properties that do not match their declared type", "event_id", event.ID, "event_key", event.EventKey, "error", err)
		}

		err := batch.Append(event.EventKey, event.UserDetails.ID, event.SentAt, event.ReceivedAt, stringProps, intProps, floatProps)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
