package model

type SimulationConfig map[string]ExperimentConfig

type VariantKey string
type EventKey string

type EventField string

type ExperimentConfig struct {
	ExperimentKey string                     `yaml:"experiment_key"`
	RandomSeed    int64                      `yaml:"random_seed"`
	Variants      map[VariantKey]VariantRole `yaml:"variants"`
	AA            ExperimentPhase            `yaml:"aa"`
	AB            ExperimentPhase            `yaml:"ab"`
	Events        map[EventKey]EventConfig   `yaml:"events"`
}

type ExperimentPhase struct {
	DurationSeconds int                             `yaml:"duration_seconds"`
	PublishAmounts  map[EventKey]map[VariantKey]int `yaml:"publish_amounts"`
}

type EventConfig struct {
	Fields map[EventField]FieldConfig `yaml:"fields"`
}

type FieldConfig struct {
	Type FieldType                        `yaml:"type"`
	AA   map[VariantKey]FieldConfigMinMax `yaml:"aa"`
	AB   map[VariantKey]FieldConfigMinMax `yaml:"ab"`
}

type FieldConfigMinMax struct {
	Min *float64 `yaml:"min"`
	Max *float64 `yaml:"max"`
}
type FieldType string

const (
	FieldTypeFloat   FieldType = "float"
	FieldTypeString  FieldType = "string"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeInteger FieldType = "integer"
)

type VariantRole string

const (
	Control   VariantRole = "control"
	Treatment VariantRole = "treatment"
)
