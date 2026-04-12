package service

import (
	"log/slog"
	"testing"
	"time"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestEnrichWithAATestDates(t *testing.T) {

	tests := []struct {
		name              string
		givenTime         time.Time
		expectedStartTime time.Time
		expectedEndTime   time.Time
	}{
		{
			name:              "given time is rounded up to the next day",
			givenTime:         time.Date(2024, 6, 1, 0, 1, 0, 0, time.UTC),
			expectedStartTime: time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
			expectedEndTime:   time.Date(2024, 6, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name:              "given time exactly at midnight still moves to next day",
			givenTime:         time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			expectedStartTime: time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
			expectedEndTime:   time.Date(2024, 6, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name:              "given time just before midnight still rounds to same day",
			givenTime:         time.Date(2024, 6, 1, 23, 59, 59, 0, time.UTC),
			expectedStartTime: time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
			expectedEndTime:   time.Date(2024, 6, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name:              "month boundary",
			givenTime:         time.Date(2024, 6, 30, 12, 0, 0, 0, time.UTC),
			expectedStartTime: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			expectedEndTime:   time.Date(2024, 7, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:              "year boundary",
			givenTime:         time.Date(2024, 12, 31, 12, 0, 0, 0, time.UTC),
			expectedStartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedEndTime:   time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
		},
	}

	service := &ExperimentService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := &experiment2.Experiment{}

			service.enrichWithAATestDates(exp, tt.givenTime)

			if !exp.AAStartTime.Equal(tt.expectedStartTime) {
				t.Errorf("Expected start time to be %v, got %v", tt.expectedStartTime, exp.AAStartTime)
			}

			if !exp.AAEndTime.Equal(tt.expectedEndTime) {
				t.Errorf("Expected end time to be %v, got %v", tt.expectedEndTime, exp.AAEndTime)
			}
		},
		)
	}

}

func TestEnrichWithExperimentStatus(t *testing.T) {

	tests := []struct {
		name           string
		currentTime    time.Time
		exp            experiment2.Experiment
		expectedStatus experiment2.ExperimentStatus
	}{
		{
			name:        "before AA start time",
			currentTime: time.Now(),
			exp: experiment2.Experiment{
				AAStartTime: time.Now().Add(24 * time.Hour),
				AAEndTime:   time.Now().Add(8 * 24 * time.Hour),
			},
			expectedStatus: experiment2.ExperimentStatusAAPlanned,
		},
		{
			name:        "during AA test",
			currentTime: time.Now(),
			exp: experiment2.Experiment{
				AAStartTime: time.Now().Add(-24 * time.Hour),
				AAEndTime:   time.Now().Add(6 * 24 * time.Hour),
			},
			expectedStatus: experiment2.ExperimentStatusAA,
		},
		{
			name:        "after AA end time with no AB dates",
			currentTime: time.Now(),
			exp: experiment2.Experiment{
				AAStartTime: time.Now().Add(-8 * 24 * time.Hour),
				AAEndTime:   time.Now().Add(-24 * time.Hour),
			},
			expectedStatus: experiment2.ExperimentStatusAAComplete,
		},
		{
			name:        "after AA end time but before AB start time",
			currentTime: time.Now(),
			exp: experiment2.Experiment{
				AAStartTime: time.Now().Add(-8 * 24 * time.Hour),
				AAEndTime:   time.Now().Add(-24 * time.Hour),
				StartTime:   pgtype.Timestamp{Valid: true, Time: time.Now().Add(24 * time.Hour)},
				EndTime:     pgtype.Timestamp{Valid: true, Time: time.Now().Add(8 * 24 * time.Hour)},
			},
			expectedStatus: experiment2.ExperimentStatusABPlanned,
		},
		{
			name:        "during AB test",
			currentTime: time.Now(),
			exp: experiment2.Experiment{
				StartTime: pgtype.Timestamp{Valid: true, Time: time.Now().Add(-4 * 24 * time.Hour)},
				EndTime:   pgtype.Timestamp{Valid: true, Time: time.Now().Add(4 * 24 * time.Hour)},
			},
			expectedStatus: experiment2.ExperimentStatusAB,
		},
		{
			name:        "after AB end time",
			currentTime: time.Now(),
			exp: experiment2.Experiment{
				StartTime: pgtype.Timestamp{Valid: true, Time: time.Now().Add(-16 * 24 * time.Hour)},
				EndTime:   pgtype.Timestamp{Valid: true, Time: time.Now().Add(-8 * 24 * time.Hour)},
			},
			expectedStatus: experiment2.ExperimentStatusComplete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ExperimentService{
				logger: slog.New(slog.NewTextHandler(nil, nil)),
			}
			service.enrichWithExperimentStatus(&tt.exp)

			if tt.exp.Status != tt.expectedStatus {
				t.Errorf("Expected status to be %v, got %v", tt.expectedStatus, tt.exp.Status)
			}
		})
	}

}
