package services

import (
	"experiment-simulator/internal/model"
	"fmt"
	"math/rand"
	"testing"
)

func TestExperimentSimulation_GetVariantToEventToAmountExcludingExposure(t *testing.T) {

	tests := []struct {
		name string
		want map[model.VariantKey]map[model.EventKey]int
		es   ExperimentSimulation
	}{
		{
			name: "aa phase with control and treatment",
			want: map[model.VariantKey]map[model.EventKey]int{
				"button_blue": map[model.EventKey]int{
					"purchase": 5000,
				},
				"button_green": map[model.EventKey]int{
					"purchase": 5000,
				},
			},
			es: ExperimentSimulation{
				ExperimentConfig: model.ExperimentConfig{
					ExperimentKey: "best_experiment_ever",
					RandomSeed:    21,
					Variants: map[model.VariantKey]model.VariantRole{
						"button_blue":  model.Control,
						"button_green": model.Treatment,
					},
					AA: model.ExperimentPhase{
						DurationSeconds: 10,
						PublishAmounts: map[model.EventKey]map[model.VariantKey]int{
							"purchase": {
								"button_blue":  5000,
								"button_green": 5000,
							},
						},
					},
					Events: map[model.EventKey]model.EventConfig{
						"purchase": {},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.es.GetVariantToEventToAmountExcludingExposure(model.ExperimentPhaseAA)
			gotLen := len(got)
			wantLen := len(tt.want)

			if gotLen != wantLen {
				t.Errorf("Got len %d , want len %d", gotLen, wantLen)
			}

			for vKey, _ := range tt.es.ExperimentConfig.Variants {
				if tt.es.ExperimentConfig.Variants[vKey] == model.Treatment {
					continue
				}

				for eKey, _ := range tt.es.ExperimentConfig.Events {
					wantCount := tt.want[vKey][eKey]
					gotCount := got[vKey][eKey]
					if tt.want[vKey][eKey] != got[vKey][eKey] {
						t.Errorf("Got %d events, want %d events", gotCount, wantCount)
					}

				}

			}

		})
	}
}

func TestDeriveSeed_Deterministic(t *testing.T) {
	seed1 := DeriveSeed(42, "aa", "purchase", "order_total")
	seed2 := DeriveSeed(42, "aa", "purchase", "order_total")

	if seed1 != seed2 {
		t.Errorf("expected same seed on repeated calls, got %d and %d", seed1, seed2)
	}
}

func TestDeriveSeed_SameRNGSequence(t *testing.T) {
	seed := DeriveSeed(42, "ab", "purchase", "order_total")

	rng1 := rand.New(rand.NewSource(seed))
	rng2 := rand.New(rand.NewSource(seed))

	for i := range 10 {
		v1 := rng1.Float64()
		v2 := rng2.Float64()
		if v1 != v2 {
			t.Errorf("value %d: expected %f == %f", i, v1, v2)
		}
	}
}

func TestDeriveSeed_DifferentParentsProduceDifferentSeeds(t *testing.T) {
	seed1 := DeriveSeed(42, "aa")
	seed2 := DeriveSeed(99, "aa")

	if seed1 == seed2 {
		t.Errorf("expected different seeds for different parents, got %d for both", seed1)
	}
}

func TestGenerateDataForField_SameFloatSequenceGivenRNG(t *testing.T) {
	seed := DeriveSeed(42, "seed1")
	sameSeed := DeriveSeed(42, "seed1")

	mnFloat64 := float64(0)
	mxFloat64 := float64(500)

	fieldType := model.FieldTypeFloat
	config := model.FieldConfigMinMax{Min: &mnFloat64, Max: &mxFloat64}

	for i := range 10000 {
		v1 := GenerateDataForField(fieldType, config, rand.NewSource(seed))
		v2 := GenerateDataForField(fieldType, config, rand.NewSource(sameSeed))

		if v1 != v2 {
			t.Errorf("value %d: expected %f == %f", i, v1, v2)
		}
	}
}

type MockUserIdService struct {
	RandSource rand.Source
}

func (m MockUserIdService) GetXUserIdsWithinExperimentAndVariant(count int, experimentKey string, wantVariantKey model.VariantKey) ([]string, error) {
	r := rand.New(m.RandSource)

	userIds := make([]string, count)

	for i := range count {
		userIds[i] = fmt.Sprintf("%d", r.Int())
	}

	return userIds, nil
}

func newAATestSimulation() (*ExperimentSimulation, func() *[]model.ExperimentParticipant) {
	mnOrderTotal := float64(1)
	mxOrderTotal := float64(500)

	es := &ExperimentSimulation{
		UserIdService: MockUserIdService{RandSource: rand.NewSource(21)},
		ExperimentConfig: model.ExperimentConfig{
			ExperimentKey: "best_experiment_ever",
			RandomSeed:    21,
			Variants: map[model.VariantKey]model.VariantRole{
				"button_blue":  model.Control,
				"button_green": model.Treatment,
			},
			AA: model.ExperimentPhase{
				DurationSeconds: 10,
				PublishAmounts: map[model.EventKey]map[model.VariantKey]int{
					"experiment_exposure": {
						"button_blue":  5000,
						"button_green": 5000,
					},
					"purchase": {
						"button_blue":  3000,
						"button_green": 3000,
					},
				},
			},
			Events: map[model.EventKey]model.EventConfig{
				"purchase": {
					Fields: map[model.EventField]model.FieldConfig{
						"order_total": {
							Type: model.FieldTypeFloat,
							AA: map[model.VariantKey]model.FieldConfigMinMax{
								"button_blue":  {Min: &mnOrderTotal, Max: &mxOrderTotal},
								"button_green": {Min: &mnOrderTotal, Max: &mxOrderTotal},
							},
						},
					},
				},
			},
		},
	}

	return es, es.GetAATestParticipantsWithActions
}

func TestExperimentSimulation_GetAATestParticipantsWithActions_TotalParticipantsMatchExposureEvents(t *testing.T) {
	es, getParticipants := newAATestSimulation()
	participants := *getParticipants()

	wantTotal := es.ExperimentConfig.AA.PublishAmounts["experiment_exposure"]["button_blue"] +
		es.ExperimentConfig.AA.PublishAmounts["experiment_exposure"]["button_green"]

	if len(participants) != wantTotal {
		t.Errorf("got %d participants, want %d", len(participants), wantTotal)
	}
}

func TestExperimentSimulation_GetAATestParticipantsWithActions_ControlPurchaseEventsMapOverExposureEvents(t *testing.T) {
	es, getParticipants := newAATestSimulation()
	participants := *getParticipants()

	wantControlPurchases := es.ExperimentConfig.AA.PublishAmounts["purchase"]["button_blue"]

	for i := range wantControlPurchases {
		if participants[i].Actions[0].EventKey != "experiment_exposure" {
			t.Errorf("control participant %d: expected experiment_exposure as 1st event but got %s", i, participants[i].Actions[0].EventKey)
		}
		if participants[i].Actions[1].EventKey != "purchase" {
			t.Errorf("control participant %d: expected purchase as 2nd event but got %s", i, participants[i].Actions[1].EventKey)
		}
	}
}

func TestExperimentSimulation_GetAATestParticipantsWithActions_TreatmentPurchaseEventsMapOverExposureEvents(t *testing.T) {
	es, getParticipants := newAATestSimulation()
	participants := *getParticipants()

	controlExposures := es.ExperimentConfig.AA.PublishAmounts["experiment_exposure"]["button_blue"]
	wantTreatmentPurchases := es.ExperimentConfig.AA.PublishAmounts["purchase"]["button_green"]

	for i := range wantTreatmentPurchases {
		idx := controlExposures + i
		if participants[idx].Actions[0].EventKey != "experiment_exposure" {
			t.Errorf("treatment participant %d: expected experiment_exposure as 1st event but got %s", idx, participants[idx].Actions[0].EventKey)
		}
		if participants[idx].Actions[1].EventKey != "purchase" {
			t.Errorf("treatment participant %d: expected purchase as 2nd event but got %s", idx, participants[idx].Actions[1].EventKey)
		}
	}
}
