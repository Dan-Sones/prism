package model

import (
	"bytes"
	"encoding/json"
)

type OutboundEventField struct {
	Value    any                        `json:"value"`
	DataType OutboundEventFieldDataType `json:"dataType"`
}

// UnmarshalJSON decodes "value" with UseNumber so JSON numbers arrive as json.Number
// rather than float64. A float64 cannot hold every int64 (anything above 2^53 loses
// precision), so the typed getters on DownstreamEvent need the original literal to
// convert int properties exactly.
func (f *OutboundEventField) UnmarshalJSON(data []byte) error {
	var raw struct {
		Value    json.RawMessage            `json:"value"`
		DataType OutboundEventFieldDataType `json:"dataType"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var value any
	if len(raw.Value) > 0 {
		dec := json.NewDecoder(bytes.NewReader(raw.Value))
		dec.UseNumber()
		if err := dec.Decode(&value); err != nil {
			return err
		}
	}

	f.Value = value
	f.DataType = raw.DataType
	return nil
}
