package services

import (
	"context"
	"data-cooking-service/internal/clients"
	"data-cooking-service/internal/repository"
	"encoding/json"
	"sync"

	"github.com/Dan-Sones/prismdbmodels/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/syncmap"
)

type Assignments map[string]map[string]string

type MicroBatchProcessorImp struct {
	cookedEventsRepository          repository.CookedEventsRepository
	assignmentClient                clients.AssignmentClient
	experimentationAssignmentClient clients.ExperimentationClient
}

func NewMicroBatchProcessorImp(repository repository.CookedEventsRepository, assignmentClient clients.AssignmentClient, experimentationAssignmentClient clients.ExperimentationClient) *MicroBatchProcessorImp {
	return &MicroBatchProcessorImp{
		cookedEventsRepository:          repository,
		assignmentClient:                assignmentClient,
		experimentationAssignmentClient: experimentationAssignmentClient,
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

	cookedEvents, err := p.cookEvents(events, assignments)
	if err != nil {
		return err
	}

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

func (p *MicroBatchProcessorImp) cookEvents(events []model.DownstreamEvent, assignments Assignments) ([]*model.CookedDownstreamEvent, error) {
	var cookedEvents []*model.CookedDownstreamEvent

	// Given the user id within the event, lookup their assignments AT the point of time in the event.
	// There could be an unlimited backlog in kafka, so we need the assignment at the point of the event, not at the point of processing.

	// Given the user id within the event, look up their assignment. Make sure to put requestor as data-cooking-service so we don't get the a/a override
	// From this will we get a list of experiment keys for the userId
	// for each experiment key, we can then look up each of those and get back the enriched experiment.
	// if the event_key is in use for that experiment then we will write a row for the event for that experiment
	// If the experiment is in the a/a period, we set is_aa to true

	// need to consider batching here. - This is LOTS of grpc requests. Caching might be needed.

	uniqueUserIds := p.deduplicateUserIds(events)

	// we need to consider that the experiment config may have been at a different time for each event, so we need to fetch the assignments per bucket per event, not globally

	userIdToBucket := make(map[string]int, len(uniqueUserIds))

	for _, event := range events {
		eventCtx := context.Background()
		bucket, ok := userIdToBucket[event.UserDetails.ID]
		if !ok {
			var err error
			bucket, err = p.assignmentClient.GetBucketForUserId(eventCtx, event.UserDetails.ID)
			if err != nil {
				return nil, err
			}
			userIdToBucket[event.UserDetails.ID] = bucket
		}

		assigmentForBucketAtEventTime, err := p.experimentationAssignmentClient.GetExperimentsAndVariantsForBucketAtTime(eventCtx, bucket, "data-cooking-service", event.SentAt)
		if err != nil {
			return nil, err
		}

		for _, exp := range assigmentForBucketAtEventTime {
			variantKeyWithinExperiment, err := p.assignmentClient.GetVariantForUserFromExperimentDetails(eventCtx, event.UserDetails.ID, exp)
			if err != nil {
				return nil, err
			}

		}

	}

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

func (p *MicroBatchProcessorImp) deduplicateUserIds(events []model.DownstreamEvent) []string {
	seen := make(map[string]struct{}, len(events))
	userIds := make([]string, 0, len(events))

	for _, event := range events {
		if _, ok := seen[event.UserDetails.ID]; ok {
			continue
		}
		seen[event.UserDetails.ID] = struct{}{}
		userIds = append(userIds, event.UserDetails.ID)
	}

	return userIds
}

//
//func (p *MicroBatchProcessorImp) getBucketsForUserIds(ctx context.Context, userIds []string) (map[string]int, error) {
//	var (
//		mu      sync.Mutex
//		buckets = make(map[string]int, len(userIds))
//	)
//
//	g, ctx := errgroup.WithContext(ctx)
//	g.SetLimit(10)
//
//	for _, userId := range userIds {
//		uid := userId
//		g.Go(func() error {
//			bucket, err := p.assignmentClient.GetBucketForUserId(ctx, uid)
//			if err != nil {
//				return err
//			}
//			mu.Lock()
//			buckets[uid] = bucket
//			mu.Unlock()
//			return nil
//		})
//	}
//
//	if err := g.Wait(); err != nil {
//		return nil, err
//	}
//	return buckets, nil
//}

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
