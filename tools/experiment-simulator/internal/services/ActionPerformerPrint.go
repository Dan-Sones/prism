package services

import (
	"experiment-simulator/internal/model"
	"fmt"
)

type ActionPerformerPrint struct{}

func (a ActionPerformerPrint) PerformAction(request model.EventRequest) {
	fmt.Printf("Performing action: %s for user: %s with properties: %v at %s\n",
		request.EventKey, request.UserDetails.Id, request.Properties, request.SentAt.Format("2006-01-02 15:04:05"))
}
