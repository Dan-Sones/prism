package errors

type FailedToCreateExperimentError struct {
	ExperimentName string
}

func (e *FailedToCreateExperimentError) Error() string {
	return "Failed to Create Experiment"
}
