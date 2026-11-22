package services

import (
	"reflect"
	"testing"
)

func TestUnMarshalMicrobatch_IntPropertiesSurviveJSONDecode(t *testing.T) {
	p := &MicroBatchProcessorImp{}

	msg := []byte(`{"id":"e1","event_key":"purchase","user_details":{"id":"u1"},"sent_at":"2026-11-22T10:00:00Z","received_at":"2026-11-22T10:00:01Z","properties":{` +
		`"quantity":{"dataType":"int","value":123},` +
		`"big":{"dataType":"int","value":9223372036854775807},` +
		`"fractional":{"dataType":"int","value":1.5},` +
		`"stringy":{"dataType":"int","value":"123"}}}`)

	events, err := p.unMarshalMicrobatch([][]byte{msg})
	if err != nil {
		t.Fatalf("unMarshalMicrobatch: %v", err)
	}

	ints, err := events[0].GetIntProperties()
	if want := map[string]int64{"quantity": 123, "big": 9223372036854775807}; !reflect.DeepEqual(ints, want) {
		t.Errorf("int properties = %v, want %v", ints, want)
	}
	if err == nil {
		t.Error("expected conversion error for fractional and stringy")
	}
}
