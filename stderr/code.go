package stderr

import (
	"encoding/json"
	"regexp"
)

var (
	codeFormat = regexp.MustCompile(`^[A-Za-z]\w*$`)
)

// Code returns a new code typed error. When placed in a chain of errors, this type of error often
// suggests the appropriate error code to be displayed in the API response.
//
// The code supplied to the function must meet the code format. By default, it requires a string that
// starts with lowercase or uppercase character, and contains only alphanumeric characters and
// underscore. The regular expression ^[A-Za-z]\w*$ is used to validate it. The default error code
// format can be changed by SetErrorCodeFormat.
func Code(code string) Error {
	if len(code) == 0 || !codeFormat.MatchString(code) {
		panic("code does not match error code format")
	}
	return &CodeError{code: code}
}

// SetErrorCodeFormat sets the global error format which guards the Code constructor.
func SetErrorCodeFormat(format string) {
	codeFormat = regexp.MustCompile(format)
}

type CodeError struct {
	code string
	next Error
}

func (e *CodeError) Code() string {
	return e.code
}

func (e *CodeError) Error() string {
	return e.code
}

func (e *CodeError) Unwrap() error {
	return e.next
}

func (e *CodeError) Is(target error) bool {
	switch ce := target.(type) {
	case *CodeError:
		return ce.code == e.code
	default:
		return false
	}
}

func (e *CodeError) wrap(err Error) {
	e.next = err
}

func (e *CodeError) asNode() (*node, error) {
	jsonBytes, err := json.Marshal(codeErrorJSON{Code: e.code})
	if err != nil {
		return nil, err
	}

	return &node{
		Type: typeCode,
		Data: jsonBytes,
	}, nil
}

func (e *CodeError) UnmarshalJSON(bytes []byte) error {
	var temp codeErrorJSON
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	e.code = temp.Code

	return nil
}

type codeErrorJSON struct {
	Code string `json:"code"`
}
