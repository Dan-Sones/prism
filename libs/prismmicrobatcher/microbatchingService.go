package services

import (
	"context"
	"log/slog"
	"time"
)

type MicroBatchingService struct {
	microBatchSize      int
	flushTimeout   time.Duration
	eventReader         EventReader
	microBatchProcessor MicrobatchProcessor
	logger              *slog.Logger
}

func NewMicroBatchingService(microBatchSize int, flushTimeout time.Duration, eventReader EventReader, microbatchProcessor MicrobatchProcessor, logger *slog.Logger) *MicroBatchingService {
	return &MicroBatchingService{
		microBatchSize:      microBatchSize,
		flushTimeout:   flushTimeout,
		eventReader:         eventReader,
		microBatchProcessor: microbatchProcessor,
		logger:              logger,
	}
}

func (m *MicroBatchingService) Start(ctx context.Context) {
	currentBatch := make([][]byte, 0, m.microBatchSize)
	m.logger.Info("Micro Batching started with microbatch size", "size", m.microBatchSize)

	batchFlushTimeout := time.NewTicker(m.flushTimeout)

	messageCh := make(chan [][]byte)

	go func() {
		for {
			messages, err := m.eventReader.PollEvents(ctx)
			if err != nil {
				m.logger.Error("Error polling events", "error", err)
				continue
			}

			if ctx.Err() != nil {
				return
			}

			if len(messages) > 0 {
				messageCh <- messages
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			remaining := m.flushFullBatches(ctx, currentBatch)
			if len(remaining) > 0 {
				err := m.microBatchProcessor.ProcessMicrobatch(ctx, remaining)
				if err != nil {
					m.logger.Error("Error processing final microbatch", "error", err)
				}
			}
			return
		case <-batchFlushTimeout.C:
			if len(currentBatch) > 0 {
				m.logger.Info("Flush timeout reached, flushing partial batch", "size", len(currentBatch))
				err := m.microBatchProcessor.ProcessMicrobatch(ctx, currentBatch)
				if err != nil {
					m.logger.Error("Error processing timeout flush", "error", err)
				} else {
					currentBatch = make([][]byte, 0, m.microBatchSize)
				}
			}
			batchFlushTimeout.Reset(m.flushTimeout)
		case messages := <-messageCh:
			for _, msg := range messages {
				currentBatch = append(currentBatch, msg)
			}

			currentBatch = m.flushFullBatches(ctx, currentBatch)
			batchFlushTimeout.Reset(m.flushTimeout)
		}
	}
}

func (m *MicroBatchingService) flushFullBatches(ctx context.Context, batch [][]byte) [][]byte {
	for len(batch) >= m.microBatchSize {
		// TODO: what if the insert fails, do we set it aside in a queue or do we not letter the buffer empty?
		// this might cause backpressure tho
		err := m.microBatchProcessor.ProcessMicrobatch(ctx, batch[:m.microBatchSize])
		if err != nil {
			m.logger.Error("Error processing microbatch", "error", err)
			return batch
		}
		batch = batch[m.microBatchSize:]
	}
	return batch
}
