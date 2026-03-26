package model

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type ParticipantEventParameters struct {
	EventField string
	properties map[string]interface{}
}

func NewParticipantEventParameters(eventField string, properties map[string]interface{}) ParticipantEventParameters {
	return ParticipantEventParameters{
		EventField: eventField,
		properties: properties,
	}
}

type ExperimentParticipant struct {
	UserId          int
	VariantKey      string
	FeatureFlagKey  string
	Actions         []ParticipantEventParameters
	ActionPerformer ActionPerformer
}

func NewExperimentParticipant(userId int, variantKey string, featureFlagKey string, performer ActionPerformer) ExperimentParticipant {
	return ExperimentParticipant{
		UserId:          userId,
		VariantKey:      variantKey,
		FeatureFlagKey:  featureFlagKey,
		ActionPerformer: performer,
	}
}

func (ep *ExperimentParticipant) AddAction(action ParticipantEventParameters) {
	ep.Actions = append(ep.Actions, action)
}

func (ep *ExperimentParticipant) PerformActionsWithDelay() {
	for _, action := range ep.Actions {
		// Simulate some delay between actions
		time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond)

		eventReq := EventRequest{
			EventKey: action.EventField,
			UserDetails: UserDetails{
				Id: strconv.Itoa(ep.UserId),
			},
			SentAt:     time.Now(),
			Properties: action.properties,
		}

		ep.ActionPerformer.PerformAction(eventReq)
	}
}
