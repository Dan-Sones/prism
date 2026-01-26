package controller

import (
	"assignment-service/internal/model/utility"
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

func WriteInternalServerError(w http.ResponseWriter) {
	problemDetail := utility.ProblemDetail{
		Title:     "Internal Server Error",
		Status:    500,
		Detail:    "An unexpected error occurred while processing your request",
		ToDisplay: "Something went wrong. Please try again later.",
	}
	WriteResponse(w, http.StatusInternalServerError, problemDetail)
}

func WriteEmptyBodyError(w http.ResponseWriter) {
	problemDetail := utility.ProblemDetail{
		Title:     "Request body is required",
		Status:    http.StatusBadRequest,
		ToDisplay: "Request body is required",
	}
	WriteResponse(w, http.StatusBadRequest, problemDetail)
}
