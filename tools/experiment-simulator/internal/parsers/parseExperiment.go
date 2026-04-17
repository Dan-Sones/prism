package parsers

import (
	"experiment-simulator/internal/model"
	"fmt"

	"github.com/goccy/go-yaml"
)

func ParseExperimentConfig(raw []byte) model.SimulationConfig {

	var config model.SimulationConfig

	err := yaml.Unmarshal(raw, &config)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	return config
}
