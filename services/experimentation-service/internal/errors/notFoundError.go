package errors

import "fmt"

type NotFoundError struct {
	Resource string
	ID       int32
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %v not found", e.Resource, e.ID)
}
