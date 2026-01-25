package utility

import "encoding/json"

type ProblemDetail struct {
	Title            string          `json:"title,omitempty"`
	Status           int             `json:"status,omitempty"`
	Detail           string          `json:"detail,omitempty"`
	ToDisplay        string          `json:"toDisplay,omitempty"`
	StructuredDetail json.RawMessage `json:"structuredDetail,omitempty"`
}
