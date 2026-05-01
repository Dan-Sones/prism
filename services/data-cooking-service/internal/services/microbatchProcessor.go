package services

import (
	"context"
	"data-cooking-service/internal/clients"
	"data-cooking-service/internal/repository"
	"encoding/json"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/Dan-Sones/prismdbmodels/model/experiment"
)

type Assignments map[string]map[string]string

type MicroBatchProcessorImp struct {
	cookedEventsRepository          repository.CookedEventsRepository
	assignmentClient                clients.AssignmentClient
	experimentClient                clients.ExperimentationExperimentClient
	experimentationAssignmentClient clients.ExperimentationAssignmentClient
}

func NewMicroBatchProcessorImp(repository repository.CookedEventsRepository, assignmentClient clients.AssignmentClient, experimentClient clients.ExperimentationExperimentClient, experimentationAssignmentClient clients.ExperimentationAssignmentClient) *MicroBatchProcessorImp {
	return &MicroBatchProcessorImp{
		cookedEventsRepository:          repository,
		assignmentClient:                assignmentClient,
		experimentClient:                experimentClient,
		experimentationAssignmentClient: experimentationAssignmentClient,
	}
}

func (p *MicroBatchProcessorImp) ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error {
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

	uniqueUserIds := p.deduplicateUserIds(events)
	userIdToBucket := make(map[string]int, len(uniqueUserIds))

	//ExpKey -> ExpDetails
	enrichedExperimentDetails := make(map[string]experiment.EnrichedExperiment)

	for _, event := range events {
		// TODO: maybe a better context type here?
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

			var experimentDetails experiment.EnrichedExperiment
			cachedExp, ok := enrichedExperimentDetails[exp.ExperimentKey]
			if ok {
				experimentDetails = cachedExp
			} else {
				experimentDetails, err = p.experimentClient.GetEnrichedExperimentByKey(eventCtx, exp.ExperimentKey)
				if err != nil {
					return nil, err
				}
			}

			eventKeyInExperiment := p.isEventKeyInExperiment(event.EventKey, experimentDetails)
			if !eventKeyInExperiment {
				continue
			} else {

				if event.SentAt.After(experimentDetails.AAStartTime) && event.SentAt.Before(experimentDetails.AAEndTime) {

				}

				cookedEvents = append(cookedEvents, &model.CookedDownstreamEvent{
					DownstreamEvent: event,
					ExperimentKey:   exp.ExperimentKey,
					VariantKey:      variantKeyWithinExperiment,
					IsAA:            event.SentAt.After(experimentDetails.AAStartTime) && event.SentAt.Before(experimentDetails.AAEndTime),
				})

			}
		}

	}

	return nil, nil
}

func (p *MicroBatchProcessorImp) isEventKeyInExperiment(eventKey string, experiment experiment.EnrichedExperiment) bool {
	for _, metric := range experiment.Metrics {
		for _, metricComponents := range metric.MetricComponents {
			if metricComponents.EventType.EventKey == eventKey {
				return true
			}
		}
	}
	return false
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
