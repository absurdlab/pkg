package check

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// String is the entrypoint for checking a string.
func String(value string) *StringCheck {
	return &StringCheck{value: value, skip: false}
}

type StringCheck struct {
	value string
	skip  bool
	err   error
}

// Required checks if the string is not empty. It registers error when the string is empty. This step cannot be skipped.
func (s *StringCheck) Required() *StringCheck {
	if len(s.value) == 0 {
		s.err = errors.New("require string is empty")
	}
	return s
}

// Optional mark the succeeding checks as optional, and skips them if the value is empty. This step cannot be skipped.
func (s *StringCheck) Optional() *StringCheck {
	return s.OptionalWhen(func(value string) bool {
		return len(value) == 0
	})
}

// OptionalWhen mark the succeeding checks as optional, and skips them if the value meets the condition.
// This step cannot be skipped.
func (s *StringCheck) OptionalWhen(condition func(value string) bool) *StringCheck {
	if condition(s.value) {
		s.skip = true
	}
	return s
}

// Equals checks if the string equals target value
func (s *StringCheck) Equals(value string) *StringCheck {
	if s.shouldSkip() {
		return s
	}

	if s.value != value {
		s.err = errors.New("string does not equal target value")
	}

	return s
}

// EqualsCaseInsensitive checks if the string equals target value, comparison is done in case-insensitive way.
func (s *StringCheck) EqualsCaseInsensitive(value string) *StringCheck {
	if s.shouldSkip() {
		return s
	}

	if strings.ToLower(s.value) != strings.ToLower(value) {
		s.err = errors.New("string does not equal target value")
	}

	return s
}

// Enum checks if the string is one of the domain value.
func (s *StringCheck) Enum(domain ...string) *StringCheck {
	if s.shouldSkip() {
		return s
	}

	for _, d := range domain {
		if d == s.value {
			return s
		}
	}

	s.err = errors.New("enum does not contain string")

	return s
}

// Regex checks string matches the regular expression.
func (s *StringCheck) Regex(regex *regexp.Regexp) *StringCheck {
	if s.shouldSkip() {
		return s
	}

	match := regex.MatchString(s.value)
	if !match {
		s.err = errors.New("string does not match regex")
	}

	return s
}

// IsEmail checks if a string is a valid email address.
func (s *StringCheck) IsEmail() *StringCheck {
	if s.shouldSkip() {
		return s
	}

	if len(s.value) < 3 && len(s.value) > 254 {
		s.err = errors.New("string is not email: invalid length")
		return s
	}

	return s.Regex(emailRegex)
}

// IsURL checks if the string satisfies the url check.
func (s *StringCheck) IsURL(f func(value *url.URL) *URLCheck) *StringCheck {
	if s.shouldSkip() {
		return s
	}

	u, err := url.Parse(s.value)
	if err != nil {
		s.err = errors.New("string is not a valid url")
		return s
	}

	s.err = f(u).Error()

	return s
}

// IsValidJSON checks if the string's content is a valid JSON.
func (s *StringCheck) IsValidJSON() *StringCheck {
	if s.shouldSkip() {
		return s
	}

	if !json.Valid([]byte(s.value)) {
		s.err = errors.New("string does not contain value JSON")
	}

	return s
}

// That performs the caller defined checks.
func (s *StringCheck) That(customCheck func(value string) error) *StringCheck {
	if s.shouldSkip() {
		return s
	}

	s.err = customCheck(s.value)

	return s
}

// Error terminates the check chain and returns any existing error.
func (s *StringCheck) Error() error {
	return s.err
}

func (s *StringCheck) shouldSkip() bool {
	return s.skip || s.err != nil
}

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)
