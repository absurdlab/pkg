package stderr

import "encoding/json"

// Message returns a message typed error. When placed in a chain of errors, this type of error
// suggests the appropriate human-readable message to include in the API response. The supplied
// message must not be empty.
func Message(message string) Error {
	if len(message) == 0 {
		panic("message is required")
	}
	return &messageError{message: message}
}

type messageError struct {
	message string
	next    Error
}

func (e *messageError) Message() string {
	return e.message
}

func (e *messageError) Is(_ error) bool {
	return false
}

func (e *messageError) wrap(err Error) {
	e.next = err
}

func (e *messageError) Unwrap() error {
	return e.next
}

func (e *messageError) Error() string {
	return e.message
}

func (e *messageError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(messageErrorJSON{Message: e.message})
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeMessage,
		Data: jsonBytes,
	}, nil
}

func (e *messageError) UnmarshalJSON(bytes []byte) error {
	var temp messageErrorJSON
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	e.message = temp.Message

	return nil
}

type messageErrorJSON struct {
	Message string `json:"message"`
}
