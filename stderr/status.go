package stderr

import (
	"encoding/json"
	"fmt"
)

// Status returns a new status typed error. When placed in a chain of errors, this type of error
// often suggests the appropriate HTTP response status.
func Status(status int) *StatusError {
	if status < 0 {
		panic("status must be non-negative")
	}
	return &StatusError{status: status}
}

type StatusError struct {
	status int
	next   Error
}

func (e *StatusError) Status() int {
	return e.status
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status: %d", e.status)
}

func (e *StatusError) Unwrap() error {
	return e.next
}

func (e *StatusError) Is(target error) bool {
	switch se := target.(type) {
	case *StatusError:
		return se.status == e.status
	default:
		return false
	}
}

func (e *StatusError) wrap(err Error) {
	e.next = err
}

func (e *StatusError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(statusErrorJSON{Status: e.status})
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeStatus,
		Data: jsonBytes,
	}, nil
}

func (e *StatusError) UnmarshalJSON(bytes []byte) error {
	var temp statusErrorJSON
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	e.status = temp.Status

	return nil
}

type statusErrorJSON struct {
	Status int `json:"status"`
}
