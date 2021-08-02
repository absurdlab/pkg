package timeplus

import "time"

// Timestamp is the number of seconds since unix epoch.
type Timestamp uint64

// IsZero returns true if this Timestamp is nil or equals to 0.
func (t *Timestamp) IsZero() bool {
	return t == nil || *t == 0
}

// Ref returns a reference to the Timestamp value.
func (t Timestamp) Ref() *Timestamp {
	return &t
}

// AddDuration adds the given time.Duration to this Timestamp and returns a new Timestamp.
func (t Timestamp) AddDuration(d time.Duration) Timestamp {
	return Timestamp(t.Int64() + int64(d/time.Second))
}

// AddSecond adds the given Second to this Timestamp and returns a new Timestamp.
func (t Timestamp) AddSecond(second Second) Timestamp {
	return Timestamp(t.Int64() + second.Int64())
}

// Time converts Timestamp to time.Time.
func (t Timestamp) Time() time.Time {
	return time.Unix(t.Int64(), 0)
}

// Int64 converts Timestamp to int64 type.
func (t Timestamp) Int64() int64 {
	if t.IsZero() {
		return 0
	}
	return int64(t)
}

// UInt64 converts Timestamp to uint64 type.
func (t Timestamp) UInt64() uint64 {
	if t.IsZero() {
		return 0
	}
	return uint64(t)
}

// Before returns true if this Timestamp happened before the given Timestamp. The method
// always returns false if at least one Timestamp is zero.
func (t Timestamp) Before(s Timestamp) bool {
	if t.IsZero() || s.IsZero() {
		return false
	}
	return t.Int64() < s.Int64()
}

// After returns true if this Timestamp happened after the given Timestamp. The method
// always returns false if at least one Timestamp is zero.
func (t Timestamp) After(s Timestamp) bool {
	if t.IsZero() || s.IsZero() {
		return false
	}
	return t.Int64() > s.Int64()
}

// Equals returns true if this Timestamp happened at the same time as the given Timestamp. The method
// always returns false if at least one Timestamp is zero.
func (t Timestamp) Equals(s Timestamp) bool {
	if t.IsZero() || s.IsZero() {
		return false
	}
	return t.Int64() == s.Int64()
}
