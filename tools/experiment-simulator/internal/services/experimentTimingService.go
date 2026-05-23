package services

import (
	"bytes"
	"context"
	"encoding/json"
	"experiment-simulator/internal/repository"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type ExperimentTimingService struct {
	experimentRepository *repository.ExperimentRepository
}

func NewExperimentTimingService(experimentRepository *repository.ExperimentRepository) *ExperimentTimingService {
	return &ExperimentTimingService{
		experimentRepository: experimentRepository,
	}
}

func (et *ExperimentTimingService) MoveAAStartToNow(expId uuid.UUID) error {
	fmt.Println("Moving AA start time to now for experiment")
	return et.experimentRepository.UpdateExperimentAATimes(
		context.Background(),
		expId,
		time.Now().UTC().Add(-1*time.Second),
		time.Now().UTC().Add(7*24*time.Hour),
	)
}

func (et *ExperimentTimingService) MoveABStartToNow(expId uuid.UUID) error {
	fmt.Println("Moving AB start time to now for experiment")
	return et.experimentRepository.UpdateExperimentABTimes(
		context.Background(),
		expId,
		time.Now().UTC().Add(-1*time.Second),
		time.Now().UTC().Add(7*24*time.Hour),
	)
}

func (et *ExperimentTimingService) EndAAPhase(expId uuid.UUID) error {
	fmt.Println("Ending AA phase for experiment")
	return et.experimentRepository.UpdateExperimentAATimes(
		context.Background(),
		expId,
		time.Now().UTC().Add(-7*24*time.Hour),
		time.Now().UTC().Add(-1*time.Second),
	)
}

func (et *ExperimentTimingService) EndABPhase(expId uuid.UUID) error {
	fmt.Println("Ending AB phase for experiment")
	return et.experimentRepository.UpdateExperimentABTimes(
		context.Background(),
		expId,
		time.Now().UTC().Add(-7*24*time.Hour),
		time.Now().UTC().Add(-1*time.Second),
	)
}

func (et *ExperimentTimingService) ProgressExperimentToABPhase(expId uuid.UUID) error {
	err := et.EndAAPhase(expId)
	if err != nil {
		return fmt.Errorf("failed to end AA phase: %w", err)
	}

	expServicePort := os.Getenv("EXPERIMENTATION_SERVICE_HTTP_PORT")
	url := fmt.Sprintf("http://localhost:%s/api/experiments/%s/begin-ab", expServicePort, expId.String())

	// the request validation requires we use these midnight timestamps in order
	// we will later override this with a db write
	nextMidnight := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour).Format(time.RFC3339Nano)
	midnightInSevenDays := time.Now().UTC().Add(7 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339Nano)
	body := map[string]interface{}{
		"start_time": nextMidnight,
		"end_time":   midnightInSevenDays,
		// TODO: This will assign 100% of the buckets to the experiment - maybe make this configurable in yml
		"bucket_allocation": 100,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	r, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	err = et.MoveABStartToNow(expId)
	if err != nil {
		return fmt.Errorf("failed to move AB start time to now: %w", err)
	}

	return nil
}
