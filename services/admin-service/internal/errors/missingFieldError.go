package errors

import "fmt"

type MissingFieldError struct {
	Field string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("missing required field: %s", e.Field)
}
