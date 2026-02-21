package problems

import (
	"encoding/json"
	"net/http"
)

type ProblemDetail struct {
	Title      string      `json:"title"`
	Status     int         `json:"status"`
	Detail     string      `json:"detail,omitempty"`
	Violations []Violation `json:"violations,omitempty"`
}

type Violation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (p ProblemDetail) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}

func NewInternalServerError() ProblemDetail {
	return ProblemDetail{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: "An unexpected error occurred",
	}
}

func NewBadRequestError(detail string) ProblemDetail {
	return ProblemDetail{
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}
