package event

import "time"

type LiveEventStatistics struct {
	MissingRates           map[string]int `json:"missingRates"`
	TotalEventsPast24Hours int            `json:"totalEventsPast24Hours"`
	TotalEventsPast7Days   int            `json:"totalEventsPast7Days"`
	LastReceivedTime       time.Time        `json:"lastReceivedTime"`
}
