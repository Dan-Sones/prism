package services

import (
	"clickhouse-writer/internal/model"
	"clickhouse-writer/internal/repository"
	"context"
	"encoding/json"
)

type MicrobatchProcessor interface {
	ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error
}

type MicroBatchProcessorImp struct {
	eventsRepository repository.EventsRepository
}

func NewMicroBatchProcessorImp(eventsRepository repository.EventsRepository) *MicroBatchProcessorImp {
	return &MicroBatchProcessorImp{
		eventsRepository: eventsRepository,
	}
}

func (p *MicroBatchProcessorImp) ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error {

	events, err := p.unMarshalMicrobatch(microbatch)
	if err != nil {
		return err
	}

	return p.eventsRepository.InsertBatch(ctx, events)
}

func (p *MicroBatchProcessorImp) unMarshalMicrobatch(microbatch [][]byte) ([]model.DownstreamEvent, error) {
	events := make([]model.DownstreamEvent, 0, len(microbatch))

	for _, msg := range microbatch {
		var event model.DownstreamEvent
		err := json.Unmarshal(msg, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
