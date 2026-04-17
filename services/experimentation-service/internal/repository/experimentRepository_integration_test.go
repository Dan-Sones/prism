package repository

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"log"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/hashicorp/go-version"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type ExperimentRepositoryIntegrationTest struct {
	name                string
	experimentsToCreate []experiment2.Experiment
	bucketId            int32
	expectedNames       []string
}

type versionedFile struct {
	entry   fs.DirEntry
	version *version.Version
}

func TestExperimentRepository_GetExperimentsAndVariantsForBucket_ShouldOnlyReturnActiveExperiments(t *testing.T) {
	// Where an active experiment is defined as an experiment where the current time falls within the a/a or a/b window
	ctx := context.Background()

	dbName := "prism"
	dbUser := "tcuser"
	dbPassword := "tcPassword"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:18.1-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Printf("failed to get container host: %s", err)
		return
	}

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Printf("failed to get container port: %s", err)
		return
	}
	portWithoutSuffix := port.Port()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		portWithoutSuffix,
		dbUser,
		dbPassword,
		dbName,
	)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	PerformMigrations(&testing.T{}, dbpool)

	now := time.Now()

	tests := []ExperimentRepositoryIntegrationTest{
		{
			name:     "only returns experiments whose A/A or A/B window is currently active",
			bucketId: 21,
			experimentsToCreate: []experiment2.Experiment{
				{
					// within a/a window
					Name:          "Experiment A",
					FeatureFlagID: "flag-experiment-a",
					AAStartTime:   now.Add(-1 * time.Hour),
					AAEndTime:     now.Add(1 * time.Hour),
					Variants: []experiment2.ExperimentVariant{
						{Name: "Foo", VariantKey: "foo", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
						{Name: "Bar", VariantKey: "bar", UpperBound: 100, LowerBound: 51, VariantType: experiment2.VariantTypeTreatment},
					},
				},
				{
					// Within A/B window
					Name:          "Experiment B",
					FeatureFlagID: "flag-experiment-b",
					AAStartTime:   now.Add(-72 * time.Hour),
					AAEndTime:     now.Add(-48 * time.Hour),
					StartTime:     pgtype.Timestamp{Time: now.Add(-24 * time.Hour), Valid: true},
					EndTime:       pgtype.Timestamp{Time: now.Add(24 * time.Hour), Valid: true},
					Variants: []experiment2.ExperimentVariant{
						{Name: "Baz", VariantKey: "baz", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
						{Name: "Bash", VariantKey: "bash", UpperBound: 100, LowerBound: 51, VariantType: experiment2.VariantTypeTreatment},
					},
				},
				{
					// intermediary gap
					Name:          "Experiment C",
					FeatureFlagID: "flag-experiment-c",
					AAStartTime:   now.Add(-72 * time.Hour),
					AAEndTime:     now.Add(-48 * time.Hour),
					StartTime:     pgtype.Timestamp{Time: now.Add(24 * time.Hour), Valid: true},
					EndTime:       pgtype.Timestamp{Time: now.Add(72 * time.Hour), Valid: true},
					Variants: []experiment2.ExperimentVariant{
						{Name: "Jeff", VariantKey: "jeff", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
						{Name: "Joe", VariantKey: "joe", UpperBound: 100, LowerBound: 51, VariantType: experiment2.VariantTypeTreatment},
					},
				},
			},
			expectedNames: []string{"Experiment A", "Experiment B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commonTestPattern(t, tt.experimentsToCreate, tt.expectedNames, tt.bucketId, dbpool)
		})
	}

	dbpool.Close()

	err = testcontainers.TerminateContainer(postgresContainer)
	if err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

func commonTestPattern(t *testing.T, givenExperiments []experiment2.Experiment, expectedNames []string, bucketId int32, pgxPool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	_, err := pgxPool.Exec(ctx, "TRUNCATE prism.experiments CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate experiments: %s", err)
	}

	// Given
	expRepo := NewExperimentRepository(pgxPool)
	baRepo := NewBucketAllocationRepository(pgxPool)

	for _, exp := range givenExperiments {
		experimentId, err := expRepo.CreateNewExperiment(ctx, exp)
		if err != nil {
			t.Fatalf("failed to create experiment %q: %s", exp.Name, err)
		}

		// AA times are only set on creation, so we need to manually set the A/B times for the experiments that have them
		if exp.StartTime.Valid {
			_, err = pgxPool.Exec(ctx,
				`UPDATE prism.experiments SET start_time = $1, end_time = $2 WHERE id = $3`,
				exp.StartTime.Time, exp.EndTime.Time, *experimentId)
			if err != nil {
				t.Fatalf("failed to update A/B times for experiment %q: %s", exp.Name, err)
			}
		}

		err = baRepo.AssignBucketToExperiment(ctx, *experimentId, int(bucketId))
		if err != nil {
			t.Fatalf("failed to assign bucket to experiment %q: %s", exp.Name, err)
		}
	}

	// When
	returnedExperiments, err := expRepo.GetExperimentsAndVariantsForBucket(ctx, bucketId)
	if err != nil {
		t.Fatalf("failed to get experiments and variants for bucket: %s", err)
	}

	// Then
	if len(returnedExperiments) != len(expectedNames) {
		actualNames := make([]string, len(returnedExperiments))
		for i, e := range returnedExperiments {
			actualNames[i] = e.Name
		}
		t.Fatalf("expected %d experiments %v, got %d %v", len(expectedNames), expectedNames, len(returnedExperiments), actualNames)
	}

	returnedNames := make(map[string]bool)
	for _, exp := range returnedExperiments {
		returnedNames[exp.Name] = true
	}

	for _, name := range expectedNames {
		if !returnedNames[name] {
			t.Errorf("expected experiment %q to be returned, but it was not", name)
		}
	}
}

func PerformMigrations(t *testing.T, pgxPool *pgxpool.Pool) {
	// TODO: instead of this entire function maybe we could spin up a flyway test container to execute migrations instead?
	t.Helper()

	relativeMigrationsPath := "../../../../infrastructure/postgres-migrations"
	files, err := os.ReadDir(relativeMigrationsPath)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %s", err)
	}

	// only consider flyway migration files
	files = slices.Collect(func(yield func(entry os.DirEntry) bool) {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if !file.IsDir() && strings.HasPrefix(file.Name(), "V") && strings.HasSuffix(file.Name(), ".sql") {
				yield(file)
			}
		}
	})

	versioned := make([]versionedFile, 0, len(files))
	for _, f := range files {
		raw := strings.Split(f.Name(), "__")[0]
		withoutV := strings.TrimPrefix(raw, "V")
		v, err := version.NewVersion(withoutV)
		if err != nil {
			t.Errorf("migration %q has invalid version %q", f.Name(), raw)
		}
		versioned = append(versioned, versionedFile{entry: f, version: v})
	}

	sort.Slice(versioned, func(i, j int) bool {
		return versioned[i].version.LessThan(versioned[j].version)
	})

	for _, vf := range versioned {
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", relativeMigrationsPath, vf.entry.Name()))
		if err != nil {
			t.Fatalf("failed to read migration file %s: %s", vf.entry.Name(), err)
		}

		log.Printf("Currently executing migration %s", vf.entry.Name())

		_, err = pgxPool.Exec(context.Background(), string(content))
		if err != nil {
			t.Fatalf("failed to execute migration %s: %s", vf.entry.Name(), err)
		}
	}
}
