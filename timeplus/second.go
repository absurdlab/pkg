package timeplus

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Second represent a clock second
type Second uint64

// Int64 converts the second to an int64 type.
func (s Second) Int64() int64 {
	return int64(s)
}

// UInt64 converts the second to an uint64 type.
func (s Second) UInt64() uint64 {
	return uint64(s)
}

// Duration converts the second to time.Duration type.
func (s Second) Duration() time.Duration {
	return time.Duration(s) * time.Second
}

// Ref returns a reference to the Second value.
func (s Second) Ref() *Second {
	return &s
}

func (s *Second) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	return s.Int64(), nil
}

func (s *Second) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source uint64
	switch src.(type) {
	case int64:
		source = uint64(src.(int64))
	case int:
		source = uint64(src.(int))
	default:
		return errors.New("incompatible type for second")
	}

	*s = Second(source)

	return nil
}
