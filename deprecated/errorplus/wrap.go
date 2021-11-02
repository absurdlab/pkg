package errorplus

import (
	"errors"
	"fmt"
)

// Wrap wraps the cause and returns an error type that can be customized to include further information.
func Wrap(cause error) *wrapError {
	return &wrapError{cause: cause}
}

// WrapF creates Wrap with created error.
func WrapF(format string, args ...interface{}) *wrapError {
	return Wrap(fmt.Errorf(format, args...))
}

// Decorate creates a wrapped error with this error.
func Decorate(err error) *wrapError {
	return &wrapError{err: err}
}

// DecorateF calls Decorate with created error.
func DecorateF(format string, args ...interface{}) *wrapError {
	return Decorate(fmt.Errorf(format, args...))
}

type wrapError struct {
	err        error
	cause      error
	statusHint int
	detail     detail
}

func (e *wrapError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *wrapError) Unwrap() error {
	return e.cause
}

func (e *wrapError) Is(target error) bool {
	if e.err == nil {
		return false
	}
	return errors.Is(e.err, target)
}

func (e *wrapError) As(target interface{}) bool {
	if e.err == nil {
		return false
	}
	return errors.As(e.err, target)
}

// StatusHint sets the HTTP response status hint on the error. A status hint represents the error's opinion of its
// representation in the final response. For example, a hint of 400 indicates the error thinks it is a user input error.
// As error returns upstream, caller may possess more knowledge about the business logic context and decides to override
// the opinion by providing another hint. For example, a not found event is usually hinted as 404, but the service
// upstream may decide it is a user input error and overrides it as a 400. The closest status hint to the leaf error
// usually prevails as the final status.
func (e *wrapError) StatusHint(status int) *wrapError {
	e.statusHint = status
	return e
}

// Field adds key value pair to the detail of the error. The key "message" is reserved, if used, may be override by
// a call to Message.
func (e *wrapError) Field(key string, value interface{}) *wrapError {
	if e.detail == nil {
		e.detail = map[string]interface{}{}
	}
	e.detail[key] = value
	return e
}

// Decorate sets the core error. Errors cannot be decorated twice, otherwise this function panics.
func (e *wrapError) Decorate(err error) *wrapError {
	if e.err != nil {
		panic("error is already decorated")
	}
	e.err = err
	return e
}

// DecorateF calls Decorate with the created error.
func (e *wrapError) DecorateF(format string, args ...interface{}) *wrapError {
	return e.Decorate(fmt.Errorf(format, args...))
}

// Wrap sets the cause of this error. The cause cannot be set twice, otherwise this function panics.
func (e *wrapError) Wrap(cause error) *wrapError {
	if e.cause != nil {
		panic("error already wraps another error")
	}
	e.cause = cause
	return e
}

// WrapF calls Wrap with the created error.
func (e *wrapError) WrapF(format string, args ...interface{}) *wrapError {
	return e.Wrap(fmt.Errorf(format, args...))
}
