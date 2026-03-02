package services

import (
	"context"
	"log/slog"
)

type MicroBatchingService struct {
	microBatchSize      int
	eventReader         EventReader
	microBatchProcessor MicrobatchProcessor
	logger              *slog.Logger
}

func NewMicroBatchingService(microBatchSize int, eventReader EventReader, microbatchProcessor MicrobatchProcessor, logger *slog.Logger) *MicroBatchingService {
	return &MicroBatchingService{
		microBatchSize:      microBatchSize,
		eventReader:         eventReader,
		microBatchProcessor: microbatchProcessor,
		logger:              logger,
	}
}

func (m *MicroBatchingService) Start(ctx context.Context) {
	currentBatch := make([][]byte, 0, m.microBatchSize)

	for {
		select {
		case <-ctx.Done():
			remaining := m.flushFullBatches(currentBatch)
			if len(remaining) > 0 {
				err := m.microBatchProcessor.ProcessMicrobatch(remaining)
				if err != nil {
					m.logger.Error("Error processing final microbatch", "error", err)
				}
			}
			return
		default:
			messages, err := m.eventReader.PollEvents(ctx)
			if err != nil {
				m.logger.Error("Error polling events", "error", err)
				continue
			}

			for _, msg := range messages {
				currentBatch = append(currentBatch, msg)
			}

			currentBatch = m.flushFullBatches(currentBatch)
		}
	}
}

func (m *MicroBatchingService) flushFullBatches(batch [][]byte) [][]byte {
	for len(batch) >= m.microBatchSize {
		err := m.microBatchProcessor.ProcessMicrobatch(batch[:m.microBatchSize])
		if err != nil {
			m.logger.Error("Error processing microbatch", "error", err)
			return batch
		}
		batch = batch[m.microBatchSize:]
	}
	return batch
}
