package model

type SimulationConfig map[string]ExperimentConfig

type ExperimentConfig struct {
	DurationSeconds int                      `yaml:"duration_seconds"`
	FeatureFlagKey  string                   `yaml:"feature_flag_key"`
	VariantKeys     []string                 `yaml:"variant_keys"`
	Events          map[string]EventConfig `yaml:"events"`
}

type EventConfig struct {
	Fields                   []map[string]FieldConfig `yaml:"fields"`
	CountToPublishForVariant map[string]int           `yaml:"count_to_publish_for_variant"`
}

type FieldConfig struct {
	Type string  `yaml:"type"`
	Min  float64 `yaml:"min"`
	Max  float64 `yaml:"max"`
}
