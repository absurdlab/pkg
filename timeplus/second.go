package timeplus

import "time"

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
