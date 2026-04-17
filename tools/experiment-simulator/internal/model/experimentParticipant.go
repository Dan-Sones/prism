package model

import (
	"time"
)

type ParticipantEventParameters struct {
	EventKey   EventKey
	properties map[EventField]interface{}
}

func NewParticipantEventParameters(eventKey EventKey, properties map[EventField]interface{}) ParticipantEventParameters {
	return ParticipantEventParameters{
		EventKey:   eventKey,
		properties: properties,
	}
}

type ExperimentParticipant struct {
	UserId        string
	VariantKey    VariantKey
	ExperimentKey string
	Actions       []ParticipantEventParameters
}

func NewExperimentParticipant(userId string, variantKey VariantKey, experimentKey string) ExperimentParticipant {
	return ExperimentParticipant{
		UserId:        userId,
		VariantKey:    variantKey,
		ExperimentKey: experimentKey,
	}
}

func (ep *ExperimentParticipant) AddAction(action ParticipantEventParameters) {
	ep.Actions = append(ep.Actions, action)
}

func (ep *ExperimentParticipant) PerformActionsWithDelay(performer ActionPerformer) {
	for _, action := range ep.Actions {
		// Simulate some delay between actions
		//time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond)

		eventReq := EventRequest{
			EventKey: action.EventKey,
			ExperimentDetails: ExperimentDetails{
				ExperimentKey: ep.ExperimentKey,
				VariantKey:    string(ep.VariantKey),
			},
			UserDetails: UserDetails{Id: ep.UserId},
			SentAt:      time.Now(),
			Properties:  action.properties,
		}

		performer.PerformAction(eventReq)
	}
}
