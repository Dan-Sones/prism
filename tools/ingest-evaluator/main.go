package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"ingest-evaluator/internal/clients"
	"ingest-evaluator/internal/repository"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Test struct {
	ExperimentKey        string
	ClickhouseRepo       *repository.CookedEventsRepositoryClickhouse
	KafkaClient          *kgo.Client
	InitialOffset        int64
	InitialHighWatermark int64
	Results              []Row
}

func NewTest(expKey string, clickhouseRepo *repository.CookedEventsRepositoryClickhouse, kafkaClient *kgo.Client) *Test {
	return &Test{
		ExperimentKey:  expKey,
		ClickhouseRepo: clickhouseRepo,
		KafkaClient:    kafkaClient,
		Results:        []Row{},
	}
}

type Row struct {
	RowCount      string
	HeadPosition  string
	Timestamp     string
	HighWatermark string
}

func main() {
	_ = godotenv.Load("../../infrastructure/.env")

	clickhouse, err := clients.NewClickhouseConnection()
	if err != nil {
		fmt.Printf("Error creating clickhouse connection: %v\n", err)
		return
	}

	kafkaClient, err := clients.GetKafkaClient()
	if err != nil {
		log.Fatalf("Error creating kafka producer: %v\n", err)
	}

	clickHouseRepo := repository.NewCookedEventsRepositoryClickhouse(clickhouse)

	test := NewTest("dev_test", clickHouseRepo, kafkaClient)
	test.InitialOffset, err = test.getCurrentCommittedOffset()
	test.InitialHighWatermark, err = test.getCurrentHighWatermark()

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				test.performInspection()
			}
		}
	}()

	fmt.Println("Press Enter to stop...")
	for {
		var input string
		fmt.Scanln(&input)
		if input == "" {
			break
		}
	}

	ticker.Stop()
	done <- true

	err = test.writeResultsToCsv()
	if err != nil {
		log.Fatalf("Error writing results to csv: %v\n", err)
	}

}

func (t *Test) performInspection() {
	now := time.Now()

	offset, err := t.getCurrentCommittedOffset()
	if err != nil {
		log.Fatalf("Error getting current committed offset: %v\n", err)
	}
	relativeOffset := offset - t.InitialOffset

	highWatermark, err := t.getCurrentHighWatermark()
	if err != nil {
		log.Fatalf("Error getting current high watermark: %v\n", err)
	}
	relativeHighWatermark := highWatermark - t.InitialHighWatermark

	rowCount, err := t.ClickhouseRepo.GetCountOfEventInPhaseForExperimentKey(context.Background(), "dev_test", "a/a")
	if err != nil {
		log.Fatalf("Error getting row count from clickhouse: %v\n", err)
	}

	fmt.Printf("High Watermark: %v Row Count: %v, Offset: %v \n", relativeHighWatermark, rowCount, relativeOffset)

	t.Results = append(t.Results, Row{
		RowCount:      fmt.Sprintf("%d", rowCount),
		HeadPosition:  fmt.Sprintf("%d", relativeOffset),
		HighWatermark: fmt.Sprintf("%d", relativeHighWatermark),
		Timestamp:     now.Format("2006-01-02 15:04:05"),
	})

}

func (t *Test) getCurrentCommittedOffset() (int64, error) {
	adm := kadm.NewClient(t.KafkaClient)
	ctx := context.Background()

	offs, err := adm.FetchOffsets(ctx, os.Getenv("DATA_COOKING_SERVICE_KAFKA_CONSUMER_GROUP_ID"))
	if err != nil {
		return 0, fmt.Errorf("error fetching offsets: %v", err)
	}

	offsetForTopic, ok := offs.Lookup(os.Getenv("KAFKA_EVENTS_TOPIC"), 0)
	if !ok {
		// skip the error here or it fails on a fresh load
		fmt.Println("WARNING: No offset found for topic, returning 0")
		return 0, nil
	}

	return offsetForTopic.At, nil
}

func (t *Test) getCurrentHighWatermark() (int64, error) {
	adm := kadm.NewClient(t.KafkaClient)
	ctx := context.Background()
	topic := os.Getenv("KAFKA_EVENTS_TOPIC")

	ends, err := adm.ListEndOffsets(ctx, topic)
	if err != nil {
		return 0, fmt.Errorf("error fetching end offsets: %v", err)
	}

	o, ok := ends.Lookup(topic, 0)
	if !ok {
		return 0, fmt.Errorf("no end offset")
	}

	return o.Offset, nil
}

func (t *Test) writeResultsToCsv() error {
	fileName := fmt.Sprintf("results_%s.csv", time.Now().Format(time.RFC3339))

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"high_watermark", "row_count", "head_position", "timestamp"}); err != nil {
		return fmt.Errorf("error writing header to csv: %v", err)
	}

	for _, row := range t.Results {
		if err := csvWriter.Write([]string{row.HighWatermark, row.RowCount, row.HeadPosition, row.Timestamp}); err != nil {
			return fmt.Errorf("error writing row to csv: %v", err)
		}
	}

	return csvWriter.Error()
}
