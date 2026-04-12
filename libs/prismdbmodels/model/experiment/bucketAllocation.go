package experiment

import "github.com/google/uuid"

type BucketAllocation struct {
	ID            uuid.UUID
	Experiment_ID uuid.UUID
	BucketNumber  int
}
