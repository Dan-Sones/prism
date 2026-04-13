package services

import (
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/parsers"
	"log"
	"os"
	"path/filepath"
)

func GetSimulation() model.SimulationConfig {
	data, err := os.ReadFile(filepath.Join("resources", "experiment.yml"))
	if err != nil {
		log.Fatal(err)
	}
	return parsers.ParseExperimentConfig(data)
}
