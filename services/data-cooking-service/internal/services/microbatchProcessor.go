package services

import (
	"context"
	"data-cooking-service/internal/clients"
	"data-cooking-service/internal/repository"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismhash"
	model2 "github.com/Dan-Sones/prismhash/model"
)

type bucketDayKey struct {
	bucket int
	day    time.Time
}

type MicroBatchProcessorImp struct {
	cookedEventsRepository          repository.CookedEventsRepository
	experimentClient                clients.ExperimentationExperimentClient
	experimentationAssignmentClient clients.ExperimentationAssignmentClient
	bucketService                   *prismhash.BucketService
	logger                          *slog.Logger
}

func NewMicroBatchProcessorImp(repository repository.CookedEventsRepository,
	experimentClient clients.ExperimentationExperimentClient,
	experimentationAssignmentClient clients.ExperimentationAssignmentClient,
	bucketService *prismhash.BucketService,
	logger *slog.Logger) *MicroBatchProcessorImp {
	return &MicroBatchProcessorImp{
		cookedEventsRepository:          repository,
		experimentClient:                experimentClient,
		experimentationAssignmentClient: experimentationAssignmentClient,
		bucketService:                   bucketService,
		logger:                          logger,
	}
}

func (p *MicroBatchProcessorImp) ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error {
	p.logger.Info("processing microbatch", "size", len(microbatch))

	events, err := p.unMarshalMicrobatch(microbatch)
	if err != nil {
		p.logger.Error("failed to unmarshal microbatch", "error", err)
		return err
	}

	cookedEvents, err := p.cookEvents(events)
	if err != nil {
		p.logger.Error("failed to cook events", "error", err)
		return err
	}

	p.logger.Info("inserting cooked events", "count", len(cookedEvents))
	err = p.cookedEventsRepository.InsertBatch(ctx, cookedEvents)
	if err != nil {
		p.logger.Error("failed to insert cooked events batch", "error", err)
		return err
	}

	p.logger.Info("microbatch processed successfully", "events_in", len(events), "events_out", len(cookedEvents))
	return nil
}

func (p *MicroBatchProcessorImp) cookEvents(events []model.DownstreamEvent) ([]*model.CookedDownstreamEvent, error) {
	var cookedEvents []*model.CookedDownstreamEvent

	uniqueUserIds := p.deduplicateUserIds(events)
	userIdToBucket := make(map[string]int, len(uniqueUserIds))

	// We can cache these per day as experiments must run from utc midnight to utc midnight.
	experimentAssignmentCache := make(map[bucketDayKey][]model2.ExperimentWithVariants)

	//ExpKey -> ExpDetails
	enrichedExperimentDetails := make(map[string]experiment.EnrichedExperiment)

	for _, event := range events {
		eventCtx := context.Background()
		bucket, ok := userIdToBucket[event.UserDetails.ID]
		if !ok {
			bucket = int(p.bucketService.GetBucketFor(event.UserDetails.ID))
			userIdToBucket[event.UserDetails.ID] = bucket
		}

		bDKey := bucketDayKey{
			bucket: bucket,
			day:    event.SentAt.UTC().Truncate(24 * time.Hour),
		}

		assigmentForBucketAtEventTime, ok := experimentAssignmentCache[bDKey]
		if !ok {
			fetched, err := p.experimentationAssignmentClient.GetExperimentsAndVariantsForBucketAtTime(eventCtx, bucket, "data-cooking-service", event.SentAt)
			if err != nil {
				p.logger.Error("failed to get experiments for bucket at time", "bucket", bucket, "sent_at", event.SentAt, "error", err)
				return nil, err
			}
			experimentAssignmentCache[bDKey] = fetched
			assigmentForBucketAtEventTime = fetched
		}

		for _, exp := range assigmentForBucketAtEventTime {
			variantKeyWithinExperiment, err := prismhash.GetVariantForExperiment(exp, event.UserDetails.ID)
			if err != nil {
				p.logger.Error("failed to get variant for user from experiment details", "user_id", event.UserDetails.ID, "experiment_key", exp.ExperimentKey, "error", err)
				return nil, err
			}

			var experimentDetails experiment.EnrichedExperiment
			cachedExp, ok := enrichedExperimentDetails[exp.ExperimentKey]
			if ok {
				experimentDetails = cachedExp
			} else {
				experimentDetails, err = p.experimentClient.GetEnrichedExperimentByKey(eventCtx, exp.ExperimentKey)
				if err != nil {
					p.logger.Error("failed to get enriched experiment", "experiment_key", exp.ExperimentKey, "error", err)
					return nil, err
				}
				enrichedExperimentDetails[exp.ExperimentKey] = experimentDetails
			}

			isAA := event.SentAt.After(experimentDetails.AAStartTime) && event.SentAt.Before(experimentDetails.AAEndTime)

			if event.EventKey == "experiment_exposure" {
				p.logger.Debug("cooking exposure event", "user_id", event.UserDetails.ID, "experiment_key", exp.ExperimentKey, "variant_key", variantKeyWithinExperiment, "is_aa", isAA)
				cookedEvents = append(cookedEvents, &model.CookedDownstreamEvent{
					DownstreamEvent: event,
					ExperimentKey:   exp.ExperimentKey,
					VariantKey:      variantKeyWithinExperiment,
					IsAA:            isAA,
				})
				continue
			}

			if !p.isEventKeyInExperiment(event.EventKey, experimentDetails) {
				p.logger.Debug("event key not in experiment, skipping", "event_key", event.EventKey, "experiment_key", exp.ExperimentKey)
				continue
			}

			p.logger.Debug("cooking metric event", "user_id", event.UserDetails.ID, "event_key", event.EventKey, "experiment_key", exp.ExperimentKey, "variant_key", variantKeyWithinExperiment, "is_aa", isAA)
			cookedEvents = append(cookedEvents, &model.CookedDownstreamEvent{
				DownstreamEvent: event,
				ExperimentKey:   exp.ExperimentKey,
				VariantKey:      variantKeyWithinExperiment,
				IsAA:            isAA,
			})
		}
	}

	return cookedEvents, nil
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
