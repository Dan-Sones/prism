package services

import (
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/parsers"
	"fmt"
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

func GetUserIdsForVariant(variantKey string) model.UserIds {
	data, err := os.ReadFile(filepath.Join("resources", fmt.Sprintf("variant-%s-users.yml", variantKey)))
	if err != nil {
		log.Fatal(err)
	}
	return parsers.ParseUserIds(data)
}
