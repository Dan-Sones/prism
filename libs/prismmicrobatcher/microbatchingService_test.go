package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestMicrobatchProcessor_ProcessMicrobatch(t *testing.T) {
	const batchSize = 10000

	tests := []struct {
		name                       string
		pollSizes                  []int
		expectedBatches            int
		expectedFirstBatchMessages []string
	}{
		{
			name:                       "exact batch",
			pollSizes:                  []int{10000},
			expectedBatches:            1,
			expectedFirstBatchMessages: []string{"message_0"},
		},
		{
			name:                       "accumulates across polls",
			pollSizes:                  []int{3000, 3000, 4000},
			expectedBatches:            1,
			expectedFirstBatchMessages: []string{"message_0"},
		},
		{
			name:                       "multiple batches with remainder (remainder should be processed)",
			pollSizes:                  []int{7000, 8000, 5000, 5000},
			expectedBatches:            3,
			expectedFirstBatchMessages: []string{"message_0", "message_10000", "message_20000"},
		},
		{
			name:            "empty",
			pollSizes:       []int{},
			expectedBatches: 0,
		},
		{
			name:                       "varied poll sizes",
			pollSizes:                  []int{100, 9000, 900, 5000, 5000},
			expectedBatches:            2,
			expectedFirstBatchMessages: []string{"message_0", "message_10000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			eventReader := NewMockEventReader()
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			batchProcessor := NewStubMicrobatchProcessor()
			microBatchingService := NewMicroBatchingService(batchSize, 10*time.Second, eventReader, batchProcessor, logger)

			programPolls(t, eventReader, tt.pollSizes)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			microBatchingService.Start(ctx)

			if got := batchProcessor.GetNumberOfProcessedBatches(); got != tt.expectedBatches {
				t.Errorf("expected %d batches, got %d", tt.expectedBatches, got)
			}

			for i, expected := range tt.expectedFirstBatchMessages {
				got := batchProcessor.GetProcessedBatch(i)[0]
				if string(got) != expected {
					t.Errorf("expected first message of batch %d to be %s, got %s", i, expected, string(got))
				}
			}

		})

	}

}

func programPolls(t *testing.T, eventReader *MockEventReader, pollSizes []int) {
	t.Helper()

	msgIndex := 0
	for _, size := range pollSizes {
		messages := make([][]byte, size)
		for i := 0; i < size; i++ {
			messages[i] = []byte(fmt.Sprintf("message_%d", msgIndex))
			msgIndex++
		}
		eventReader.AddPoll(messages)
	}
}
