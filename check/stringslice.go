package check

import (
	"errors"
	"fmt"
)

// StringSlice is the entrypoint of checking a slice of strings.
func StringSlice(slice []string) *StringSliceCheck {
	return &StringSliceCheck{
		value: slice,
	}
}

type StringSliceCheck struct {
	value []string
	skip  bool
	err   error
}

// NotNil checks the slice is not nil. It cannot be skipped.
func (s *StringSliceCheck) NotNil() *StringSliceCheck {
	if s.value == nil {
		s.err = errors.New("slice is required, but is nil")
	}

	return s
}

// Optional skips the following checks if the length of the slice is 0.
func (s *StringSliceCheck) Optional() *StringSliceCheck {
	if len(s.value) == 0 {
		s.skip = true
	}

	return s
}

// NotEmpty checks the slice length is greater than 0.
func (s *StringSliceCheck) NotEmpty() *StringSliceCheck {
	if s.shouldSkip() {
		return s
	}

	if len(s.value) == 0 {
		s.err = errors.New("slice is empty")
	}

	return s
}

// Each performs checks on the slice elements. The first error breaks the check.
func (s *StringSliceCheck) Each(f func(elem string) *StringCheck) *StringSliceCheck {
	if s.shouldSkip() {
		return s
	}

	for i, elem := range s.value {
		if err := f(elem).Error(); err != nil {
			s.err = fmt.Errorf("slice element at position %d is invalid: %s", i, err)
			break
		}
	}

	return s
}

// All is an alias to Each
func (s *StringSliceCheck) All(f func(elem string) *StringCheck) *StringSliceCheck {
	return s.Each(f)
}

// Any performs checks on the slice elements, and expect at least one element to not return error.
func (s *StringSliceCheck) Any(f func(elem string) *StringCheck) *StringSliceCheck {
	if s.shouldSkip() {
		return s
	}

	for _, elem := range s.value {
		if err := f(elem).Error(); err == nil {
			return s
		}
	}

	s.err = errors.New("expect at least one slice element to meet expectation")

	return s
}

// That performs a custom check on the slice.
func (s *StringSliceCheck) That(customCheck func(slice []string) error) *StringSliceCheck {
	if s.shouldSkip() {
		return s
	}

	s.err = customCheck(s.value)

	return s
}

// Error returns the error during the checks.
func (s *StringSliceCheck) Error() error {
	return s.err
}

func (s *StringSliceCheck) shouldSkip() bool {
	return s.skip || s.err != nil
}
