package parsers

import (
	"experiment-simulator/internal/model"
	"log"

	"github.com/goccy/go-yaml"
)

func ParseUserIds(raw []byte) model.UserIds {
	var userIds model.UserIds
	err := yaml.Unmarshal(raw, &userIds)
	if err != nil {
		log.Fatal(err)
	}

	return userIds
}
