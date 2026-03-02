package services

import "clickhouse-writer/internal/model"

type MicrobatchProcessor interface {
	ProcessMicrobatch(microbatch [][]byte) error
}

type ClickhouseMicrobatchProcessor struct {
}

func NewClickhouseMicrobatchProcessor() *ClickhouseMicrobatchProcessor {
	return &ClickhouseMicrobatchProcessor{}
}

func (p *ClickhouseMicrobatchProcessor) ProcessMicrobatch(microbatch []model.DownstreamEvent) error {
	return nil
}
