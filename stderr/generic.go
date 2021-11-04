package stderr

import (
	"encoding/json"
	"errors"
)

func generic(err error) Error {
	if err == nil {
		panic("generic error must not be nil")
	}
	return &generateError{err: err}
}

type generateError struct {
	err  error
	next Error
}

func (e *generateError) Is(err error) bool {
	return errors.Is(e.err, err)
}

func (e *generateError) Error() string {
	return e.err.Error()
}

func (e *generateError) Unwrap() error {
	return e.next
}

func (e *generateError) wrap(err Error) {
	e.next = err
}

func (e *generateError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(genericErrorJSON{Error: e.Error()})
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeGeneric,
		Data: jsonBytes,
	}, nil
}

func (e *generateError) UnmarshalJSON(bytes []byte) error {
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
