package stderr

import "encoding/json"

const (
	typeStatus  = "status"
	typeCode    = "code"
	typeMessage = "message"
	typeParams  = "params"
	typeGeneric = "generic"
)

// Error is the standard interface implemented by errors provided by this package.
type Error interface {
	error
	json.Unmarshaler
	wrap(err Error)
	asNode() (*node, error)
	Unwrap() error
	Is(target error) bool
}

// Chain wraps the supplied errors in sequence and returns the first error. Errors that implement the Error
// interface are wrapped as is, other errors are normalized into a generic typed error first.
func Chain(errors ...error) error {
	if len(errors) == 0 {
		return nil
	}

	var head = normalize(errors[0])
	if len(errors) == 1 {
		return head
	}

	var cursor = head
	for i := 1; i < len(errors); i++ {
		this := normalize(errors[i])
		cursor.wrap(this)
		cursor = this
	}

	return head
}

func normalize(err error) Error {
	if err == nil {
		return nil
	}

	if stderr, ok := err.(Error); ok {
		return stderr
	}

	return generic(err)
}
