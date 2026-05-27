package experiment

type ExperimentStatus string

const (
	ExperimentStatusCancelled  ExperimentStatus = "cancelled"
	ExperimentStatusAAPlanned  ExperimentStatus = "aa-planned"
	ExperimentStatusAA         ExperimentStatus = "aa"
	ExperimentStatusAAComplete ExperimentStatus = "aa-complete"
	ExperimentStatusABPlanned  ExperimentStatus = "ab-planned"
	ExperimentStatusAB         ExperimentStatus = "ab"
	ExperimentStatusComplete   ExperimentStatus = "ab-complete"
)
