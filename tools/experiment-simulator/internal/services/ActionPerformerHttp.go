package services

import (
	"bytes"
	"encoding/json"
	"experiment-simulator/internal/model"
	"fmt"
	"log"
	"net/http"
)

type ActionPerformerHttp struct {
	Host string
	Port int
}

func NewActionPerformerHttp(host string, port int) *ActionPerformerHttp {
	return &ActionPerformerHttp{
		Host: host,
		Port: port,
	}
}

func (a *ActionPerformerHttp) PerformAction(request model.EventRequest) {

	data, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error marshalling request: %v", err)
	}

	res, err := http.Post(fmt.Sprintf("http://%s:%d/event", a.Host, a.Port), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Error response from server: %s", res.Status)
	}
}
