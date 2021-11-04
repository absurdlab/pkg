package stderr

import (
	"encoding/json"
	"fmt"
)

// Status returns a new status typed error. When placed in a chain of errors, this type of error
// often suggests the appropriate HTTP response status.
func Status(status int) Error {
	if status < 0 {
		panic("status must be non-negative")
	}
	return &statusError{status: status}
}

type statusError struct {
	status int
	next   Error
}

func (e *statusError) Status() int {
	return e.status
}

func (e *statusError) Error() string {
	return fmt.Sprintf("status: %d", e.status)
}

func (e *statusError) Unwrap() error {
	return e.next
}

func (e *statusError) Is(target error) bool {
	switch se := target.(type) {
	case *statusError:
		return se.status == e.status
	default:
		return false
	}
}

func (e *statusError) wrap(err Error) {
	e.next = err
}

func (e *statusError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(statusErrorJSON{Status: e.status})
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeStatus,
		Data: jsonBytes,
	}, nil
}

func (e *statusError) UnmarshalJSON(bytes []byte) error {
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
