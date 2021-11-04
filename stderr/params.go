package stderr

import (
	"encoding/json"
	"strings"
)

// Params returns a param typed error. When placed in a chain of errors, this type of error is used to
// provide parameterized context about the error.
func Params(keysAndValues ...interface{}) Error {
	if l := len(keysAndValues); l == 0 || l%2 != 0 {
		panic("keys and values must be in pairs")
	}

	params := map[string]interface{}{}

	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			panic("key must be a string")
		}
		params[key] = keysAndValues[i+1]
	}

	return &ParamsError{params: params}
}

type ParamsError struct {
	params map[string]interface{}
	next   Error
}

func (e *ParamsError) Params() map[string]interface{} {
	return e.params
}

func (e *ParamsError) Is(_ error) bool {
	return false
}

func (e *ParamsError) Error() string {
	if len(e.params) == 0 {
		return "params: <empty>"
	}

	var keys []string
	for k := range e.params {
		keys = append(keys, k)
	}

	return "params: " + strings.Join(keys, ", ")
}

func (e *ParamsError) Unwrap() error {
	return e.next
}

func (e *ParamsError) wrap(err Error) {
	e.next = err
}

func (e *ParamsError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(e.params)
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeParams,
		Data: jsonBytes,
	}, nil
}

func (e *ParamsError) UnmarshalJSON(bytes []byte) error {
	var temp = make(map[string]interface{})
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	e.params = temp

	return nil
}
