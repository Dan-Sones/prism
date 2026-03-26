package services

import (
	"context"
	"data-cooking-service/internal/clients"
	"data-cooking-service/internal/repository"
	"encoding/json"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type Assignments map[string]map[string]string

type MicroBatchProcessorImp struct {
	cookedEventsRepository repository.CookedEventsRepository
	assignmentClient       clients.AssignmentClient
}

func NewMicroBatchProcessorImp(repository repository.CookedEventsRepository, assignmentClient clients.AssignmentClient) *MicroBatchProcessorImp {
	return &MicroBatchProcessorImp{
		cookedEventsRepository: repository,
		assignmentClient:       assignmentClient,
	}
}

func (p *MicroBatchProcessorImp) ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error {
	// the microbatch itself is a slice of raw byte messages from kafka, we need to unmarshal these into a downstream event model
	// we then need to get the user id and store it in a list
	// we then need to make a grpc request to the assignment service to get the experiment ids and groups for each user id
	// we then need to for each experiment insert a row into clickhouse with the event data and the experiment key and variant key

	events, err := p.unMarshalMicrobatch(microbatch)
	if err != nil {
		return err
	}

	assignments, err := p.getExperimentsForUserIds(ctx, events)
	if err != nil {
		return err
	}

	cookedEvents := p.cookEvents(events, assignments)

	err = p.cookedEventsRepository.InsertBatch(ctx, cookedEvents)
	if err != nil {
		return err
	}

	return nil
}

func (p *MicroBatchProcessorImp) getExperimentsForUserIds(ctx context.Context, events []model.DownstreamEvent) (Assignments, error) {
	userIds := make([]string, 0, len(events))
	for _, event := range events {
		userIds = append(userIds, event.UserDetails.ID)
	}

	return p.assignmentClient.GetExperimentsAndVariantsForUsers(ctx, userIds)
}

func (p *MicroBatchProcessorImp) cookEvents(events []model.DownstreamEvent, assignments Assignments) []*model.CookedDownstreamEvent {
	var cookedEvents []*model.CookedDownstreamEvent

	for _, event := range events {
		for experimentKey, variantKey := range assignments[event.UserDetails.ID] {
			cookedEvents = append(cookedEvents, &model.CookedDownstreamEvent{
				DownstreamEvent: event,
				VariantKey:      variantKey,
				ExperimentKey:   experimentKey,
			})
		}
	}

	return cookedEvents
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
