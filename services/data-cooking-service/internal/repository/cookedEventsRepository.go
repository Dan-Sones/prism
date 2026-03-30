package repository

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Dan-Sones/prismdbmodels/model"
)

type CookedEventsRepository interface {
	InsertBatch(ctx context.Context, cookedEvents []*model.CookedDownstreamEvent) error
}

type CookedEventsRepositoryClickhouse struct {
	connection driver.Conn
}

func NewCookedEventsRepositoryClickhouse(connection driver.Conn) *CookedEventsRepositoryClickhouse {
	return &CookedEventsRepositoryClickhouse{
		connection: connection,
	}
}

func (r *CookedEventsRepositoryClickhouse) InsertBatch(ctx context.Context, cookedEvents []*model.CookedDownstreamEvent) error {
	insertCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	batch, err := r.connection.PrepareBatch(insertCtx, "INSERT INTO cooked_events")
	if err != nil {
		return err
	}

	for _, cookedEvent := range cookedEvents {
		err := batch.Append(cookedEvent.ExperimentKey, cookedEvent.VariantKey, cookedEvent.EventKey, cookedEvent.UserDetails.ID, cookedEvent.SentAt, cookedEvent.ReceivedAt, cookedEvent.GetStringProperties(), cookedEvent.GetIntProperties(), cookedEvent.GetFloatProperties())
		if err != nil {
			return err
		}
	}
	return batch.Send()
}
