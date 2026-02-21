package controller

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(response)
	if err != nil {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	w.Write(body)
}
