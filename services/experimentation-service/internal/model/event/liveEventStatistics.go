package event

import "time"

type LiveEventStatistics struct {
	MissingRates           map[string]int64 `json:"missingRates"`
	TotalEventsPast24Hours int64            `json:"totalEventsPast24Hours"`
	TotalEventsPast7Days   int64            `json:"totalEventsPast7Days"`
	LastReceivedTime       time.Time        `json:"lastReceivedTime"`
}
