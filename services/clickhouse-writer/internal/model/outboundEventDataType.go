package model

type OutboundEventFieldDataType string

const (
	OutboundEventFieldDataTypeString    OutboundEventFieldDataType = "string"
	OutboundEventFieldDataTypeBoolean   OutboundEventFieldDataType = "boolean"
	OutboundEventFieldDataTypeInt       OutboundEventFieldDataType = "int"
	OutboundEventFieldDataTypeFloat     OutboundEventFieldDataType = "float"
	OutboundEventFieldDataTypeTimestamp OutboundEventFieldDataType = "timestamp"
)
