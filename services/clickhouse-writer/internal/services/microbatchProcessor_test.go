package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type capturingEventsRepository struct {
	events []model.DownstreamEvent
}

func (r *capturingEventsRepository) InsertBatch(_ context.Context, events []model.DownstreamEvent) error {
	r.events = events
	return nil
}

func TestProcessMicrobatch_IntPropertiesSurviveJSONDecode(t *testing.T) {
	repo := &capturingEventsRepository{}
	processor := NewMicroBatchProcessorImp(repo)

	msg := []byte(`{"id":"e1","event_key":"purchase","user_details":{"id":"u1"},"sent_at":"2026-11-22T10:00:00Z","received_at":"2026-11-22T10:00:01Z","properties":{` +
		`"quantity":{"dataType":"int","value":123},` +
		`"big":{"dataType":"int","value":9007199254740993},` +
		`"fractional":{"dataType":"int","value":1.5},` +
		`"stringy":{"dataType":"int","value":"123"},` +
		`"price":{"dataType":"float","value":9.99}}}`)

	if err := processor.ProcessMicrobatch(context.Background(), [][]byte{msg}); err != nil {
		t.Fatalf("ProcessMicrobatch: %v", err)
	}
	if len(repo.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(repo.events))
	}

	ints, err := repo.events[0].GetIntProperties()
	if want := map[string]int64{"quantity": 123, "big": 9007199254740993}; !reflect.DeepEqual(ints, want) {
		t.Errorf("int properties = %v, want %v", ints, want)
	}
	if err == nil {
		t.Error("expected conversion error for fractional and stringy")
	}

	floats, err := repo.events[0].GetFloatProperties()
	if err != nil || !reflect.DeepEqual(floats, map[string]float64{"price": 9.99}) {
		t.Errorf("float properties = %v, err = %v", floats, err)
	}
}
