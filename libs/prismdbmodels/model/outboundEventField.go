package model

type OutboundEventField struct {
	Value    any                        `json:"value"`
	DataType OutboundEventFieldDataType `json:"dataType"`
}
