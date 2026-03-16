package services

import (
	"context"
	"encoding/json"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type MicroBatchProcessorImp struct {
}

func NewMicroBatchProcessorImp() *MicroBatchProcessorImp {
	return &MicroBatchProcessorImp{}
}

func (p *MicroBatchProcessorImp) ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error {
	// the microbatch itself is a slice of raw byte messages from kafka, we need to unmarshal these into a downstream event model
	// we then need to get the user id and store it in a list
	// we then need to make a grpc request to the assignment service to get the experiment ids and groups for each user id
	// we then need to for each experiment insert a row into clickhouse with the event data and the experiment key and variant key
	return nil
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
