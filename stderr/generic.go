package stderr

import (
	"encoding/json"
	"errors"
)

func generic(err error) *GenericError {
	if err == nil {
		panic("generic error must not be nil")
	}
	return &GenericError{err: err}
}

type GenericError struct {
	err  error
	next Error
}

func (e *GenericError) Is(err error) bool {
	return errors.Is(e.err, err)
}

func (e *GenericError) As(target interface{}) bool {
	return errors.As(e.err, target)
}

func (e *GenericError) Error() string {
	return e.err.Error()
}

func (e *GenericError) Unwrap() error {
	return e.next
}

func (e *GenericError) wrap(err Error) {
	e.next = err
}

func (e *GenericError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(genericErrorJSON{Error: e.Error()})
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeGeneric,
		Data: jsonBytes,
	}, nil
}

func (e *GenericError) UnmarshalJSON(bytes []byte) error {
	var temp genericErrorJSON
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	e.err = errors.New(temp.Error)

	return nil
}

type genericErrorJSON struct {
	Error string `json:"error"`
}
