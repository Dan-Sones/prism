package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Dan-Sones/prismdbmodels/model"
)

type CookedEventsRepository interface {
	InsertBatch(ctx context.Context, cookedEvents []*model.CookedDownstreamEvent) error
}

type CookedEventsRepositoryClickhouse struct {
	connection driver.Conn
	logger     *slog.Logger
}

func NewCookedEventsRepositoryClickhouse(connection driver.Conn, logger *slog.Logger) *CookedEventsRepositoryClickhouse {
	return &CookedEventsRepositoryClickhouse{
		connection: connection,
		logger:     logger,
	}
}

func (r *CookedEventsRepositoryClickhouse) InsertBatch(ctx context.Context, cookedEvents []*model.CookedDownstreamEvent) error {
	fmt.Printf("Inserting batch of %d cooked events\n", len(cookedEvents))
	insertCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	batch, err := r.connection.PrepareBatch(insertCtx, "INSERT INTO cooked_events")
	if err != nil {
		return err
	}

	for _, cookedEvent := range cookedEvents {
		stringProps, stringErr := cookedEvent.GetStringProperties()
		intProps, intErr := cookedEvent.GetIntProperties()
		floatProps, floatErr := cookedEvent.GetFloatProperties()
		// Drop the mistyped properties but keep the event: failing here would block the
		// whole microbatch from committing and redeliver it forever.
		if err := errors.Join(stringErr, intErr, floatErr); err != nil {
			r.logger.Warn("dropping event properties that do not match their declared type", "event_id", cookedEvent.ID, "event_key", cookedEvent.EventKey, "experiment_key", cookedEvent.ExperimentKey, "error", err)
		}

		err := batch.Append(cookedEvent.ExperimentKey, cookedEvent.VariantKey, cookedEvent.EventKey, cookedEvent.UserDetails.ID, cookedEvent.SentAt, cookedEvent.ReceivedAt, stringProps, intProps, floatProps, cookedEvent.IsAA)
		if err != nil {
			return err
		}
	}
	return batch.Send()
}
